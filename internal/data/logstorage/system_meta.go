package logstorage

import (
	"context"
	"fmt"
	"log/slog"

	"bytes"
	"encoding/json"
	"sort"

	"github.com/dgraph-io/badger/v3"
)

type MetaSystem struct {
	storage *Storage
}

func NewMetaSystem(storage *Storage) *MetaSystem {
	return &MetaSystem{storage: storage}
}

func (sys *MetaSystem) ListStreams() (map[string]StreamConfig, error) {
	streams := map[string]StreamConfig{}

	err := sys.storage.db.View(func(txn *badger.Txn) error {
		var err error
		streams, err = fetchStreamConfigs(txn)
		return err
	})

	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (sys *MetaSystem) ListStreamFields(stream string) ([]string, error) {
	fieldsMap := map[string]struct{}{}

	err := sys.storage.db.View(func(txn *badger.Txn) error {
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

func (sys *MetaSystem) GetStreamConfig(stream string) (StreamConfig, error) {
	var streamConfig StreamConfig

	err := sys.storage.db.View(func(txn *badger.Txn) error {
		var err error
		streamConfig, err = fetchStreamConfig(txn, stream)
		return err
	})

	if err != nil {
		return StreamConfig{}, err
	}

	return streamConfig, nil
}

func (sys *MetaSystem) ConfigureStream(stream string, config StreamConfig) error {
	return sys.storage.db.Update(func(txn *badger.Txn) error {
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

func (sys *MetaSystem) DeleteStream(ctx context.Context, stream string) error {
	return sys.storage.db.Update(func(txn *badger.Txn) error {
		prefixes := []string{
			fmt.Sprintf("entry:%s:", stream),
			fmt.Sprintf("index:%s:", stream),
		}

		for _, prefix := range prefixes {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(prefix)
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				key := it.Item().KeyCopy(nil)
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
		}

		streamKey := []byte(fmt.Sprintf("stream:%s", stream))
		if err := txn.Delete(streamKey); err != nil {
			return fmt.Errorf("could not delete stream config '%s': %w", stream, err)
		}

		return nil
	})
}

func getOrCreateStreamConfig(txn *badger.Txn, stream string) (StreamConfig, error) {
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

func fetchStreamConfigs(txn *badger.Txn) (map[string]StreamConfig, error) {
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

func fetchStreamConfig(txn *badger.Txn, stream string) (StreamConfig, error) {
	var streamConfig StreamConfig

	item, err := txn.Get([]byte(fmt.Sprintf("stream:%s", stream)))
	if err == badger.ErrKeyNotFound {
		return StreamConfig{}, nil
	} else if err != nil {
		return StreamConfig{}, err
	}

	err = item.Value(func(val []byte) error {
		if len(val) > 0 {
			if err := json.Unmarshal(val, &streamConfig); err != nil {
				return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
			}
		}
		return nil
	})

	if err != nil {
		return StreamConfig{}, err
	}

	return streamConfig, nil
}
