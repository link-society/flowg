package logstorage

import (
	"bytes"
	"context"
	"log/slog"

	"encoding/json"
	"fmt"

	"sort"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/app/logging"
)

type Storage struct {
	db *badger.DB
	gc *garbageCollector
}

func NewStorage(dbPath string) (*Storage, error) {
	opts := badger.
		DefaultOptions(dbPath).
		WithLogger(&logging.BadgerLogger{Channel: "logstorage"}).
		WithCompression(options.ZSTD)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	gc := newGarbageCollector(db, 5*time.Minute)
	gc.Start()

	return &Storage{db: db, gc: gc}, nil
}

func (s *Storage) Close() error {
	s.gc.Stop()
	return s.db.Close()
}

func (s *Storage) Append(
	ctx context.Context,
	stream string,
	logEntry *LogEntry,
) ([]byte, error) {
	key := logEntry.NewDbKey(stream)
	val, err := json.Marshal(logEntry)
	if err != nil {
		slog.DebugContext(
			ctx,
			"Could not marshal log entry",
			"channel", "storage",
			"stream", stream,
			"error", err.Error(),
		)
		return nil, fmt.Errorf("could not marshal log entry: %w", err)
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		slog.DebugContext(
			ctx,
			"Fetch stream configuration",
			"channel", "storage",
			"stream", stream,
		)

		streamConfig, err := s.getOrCreateStreamConfig(txn, stream)
		if err != nil {
			return err
		}

		slog.DebugContext(
			ctx,
			"Save log entry in BadgerDB",
			"channel", "storage",
			"stream", stream,
			"key", key,
		)

		entry := badger.NewEntry(key, val)
		if streamConfig.RetentionTime > 0 {
			entry = entry.WithTTL(streamConfig.RetentionTime)
		}

		if err := txn.SetEntry(entry); err != nil {
			return fmt.Errorf(
				"could not add log entry '%s' to stream '%s': %w",
				key, stream, err,
			)
		}

		for field, value := range logEntry.Fields {
			slog.DebugContext(
				ctx,
				"Save field index in BadgerDB",
				"channel", "storage",
				"stream", stream,
				"key", key,
				"field", field,
			)

			fieldIndex := newFieldIndex(txn, stream, field, value)
			if err := fieldIndex.AddKey(key); err != nil {
				return fmt.Errorf(
					"could not add field index '%s' of log entry '%s' to stream '%s': %w",
					field, key, stream, err,
				)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *Storage) Query(
	ctx context.Context,
	stream string,
	from, to time.Time,
	filter Filter,
) ([]LogEntry, error) {
	results := []LogEntry{}

	err := s.db.View(func(txn *badger.Txn) error {
		timeKeys := []string{}

		streamPrefix := []byte(fmt.Sprintf("entry:%s:", stream))
		fromPrefix := []byte(fmt.Sprintf("entry:%s:%020d:", stream, from.UnixMilli()))
		toPrefix := fmt.Sprintf("entry:%s:%020d:", stream, to.UnixMilli())

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = streamPrefix
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(fromPrefix); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.KeyCopy(nil))

			if key < toPrefix {
				timeKeys = append(timeKeys, key)
			} else {
				break
			}
		}

		var filteredKeys []string
		if filter != nil {
			fieldsIndex := newFieldsIndex(txn, stream)
			filteredKeys = fieldsIndex.Filter(filter, timeKeys)
		} else {
			filteredKeys = timeKeys
		}

		sort.Sort(sort.Reverse(sort.StringSlice(filteredKeys)))

		for _, key := range filteredKeys {
			slog.DebugContext(
				ctx,
				"Fetching log entry",
				"channel", "storage",
				"stream", stream,
				"key", key,
			)

			item, err := txn.Get([]byte(key))
			if err != nil {
				return fmt.Errorf(
					"could not fetch log entry '%s' from stream '%s': %w",
					key, stream, err,
				)
			}

			var entry LogEntry
			err = item.Value(func(val []byte) error {
				if err := json.Unmarshal(val, &entry); err != nil {
					return fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
				}

				return nil
			})
			if err != nil {
				return err
			}

			results = append(results, entry)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Storage) Purge(ctx context.Context, stream string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		prefixes := []string{
			fmt.Sprintf("entry:%s:", stream),
			fmt.Sprintf("index:%s:", stream),
		}

		keys := [][]byte{
			[]byte(fmt.Sprintf("stream:%s", stream)),
		}

		for _, prefix := range prefixes {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(prefix)
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				keys = append(keys, it.Item().KeyCopy(nil))
			}
		}

		for _, key := range keys {
			slog.DebugContext(
				ctx,
				"Purging key from BadgerDB",
				"channel", "storage",
				"stream", stream,
				"key", key,
			)

			if err := txn.Delete(key); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from stream '%s': %w",
					key, stream, err,
				)
			}
		}

		return nil
	})
}

func (s *Storage) ListStreams() (map[string]StreamConfig, error) {
	streams := map[string]StreamConfig{}

	err := s.db.View(func(txn *badger.Txn) error {
		var err error
		streams, err = fetchStreamconfigs(txn)
		return err
	})

	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (s *Storage) ListStreamFields(stream string) ([]string, error) {
	fieldsMap := map[string]struct{}{}

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := []byte(fmt.Sprintf("index:%s:field:", stream))
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = prefix
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			keyParts := bytes.SplitN(key[len(prefix):], []byte(":"), 2)
			if len(keyParts) > 0 {
				fieldName := string(keyParts[0])
				fieldsMap[fieldName] = struct{}{}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	fields := make([]string, 0, len(fieldsMap))
	for field := range fieldsMap {
		fields = append(fields, field)
	}

	sort.Strings(fields)

	return fields, nil
}

func (s *Storage) ConfigureStream(stream string, config StreamConfig) error {
	return s.db.Update(func(txn *badger.Txn) error {
		streamKey := []byte(fmt.Sprintf("stream:%s", stream))
		configVal, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
		}

		if err := txn.Set(streamKey, configVal); err != nil {
			return fmt.Errorf("could not save stream config '%s': %w", stream, err)
		}

		return nil
	})
}

func (s *Storage) getOrCreateStreamConfig(txn *badger.Txn, stream string) (StreamConfig, error) {
	var streamConfig StreamConfig

	streamKey := []byte(fmt.Sprintf("stream:%s", stream))
	switch streamConfigItem, err := txn.Get(streamKey); {
	case err != nil && err != badger.ErrKeyNotFound:
		return StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)

	case err == badger.ErrKeyNotFound:
		err := txn.Set(streamKey, []byte(""))
		if err != nil {
			return StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}

	case err == nil:
		err := streamConfigItem.Value(func(val []byte) error {
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return fmt.Errorf(
						"could not unmarshal stream config '%s': %w",
						stream, err,
					)
				}
			}
			return nil
		})
		if err != nil {
			return StreamConfig{}, err
		}
	}

	return streamConfig, nil
}

func fetchStreamconfigs(txn *badger.Txn) (map[string]StreamConfig, error) {
	streams := map[string]StreamConfig{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte("stream:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		stream := string(it.Item().Key()[7:])

		var streamConfig StreamConfig
		err := it.Item().Value(func(val []byte) error {
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		streams[stream] = streamConfig
	}

	return streams, nil
}
