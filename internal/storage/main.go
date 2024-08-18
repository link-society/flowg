package storage

import (
	"encoding/json"
	"fmt"

	"sort"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
)

type Storage struct {
	db *badger.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	opts := badger.
		DefaultOptions(dbPath).
		WithLogger(&serverLogger{}).
		WithCompression(options.ZSTD)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Append(stream string, entry *LogEntry) ([]byte, error) {
	key := entry.NewDbKey(stream)
	val, err := json.Marshal(entry)
	if err != nil {
		return nil, &MarshalError{Reason: err}
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, val); err != nil {
			return &PersistError{Operation: "append", Reason: err}
		}

		for field, value := range entry.Fields {
			fieldIndex := newFieldIndex(txn, stream, field, value)
			if err := fieldIndex.AddKey(key); err != nil {
				return &PersistError{Operation: "field-index", Reason: err}
			}
		}

		streamKey := []byte(fmt.Sprintf("stream:%s", stream))
		if err := txn.Set(streamKey, []byte{}); err != nil {
			return &PersistError{Operation: "stream-index", Reason: err}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *Storage) Query(stream string, from, to time.Time, filter Filter) ([]LogEntry, error) {
	results := []LogEntry{}

	err := s.db.View(func(txn *badger.Txn) error {
		timeKeys := []string{}

		streamPrefix := []byte(fmt.Sprintf("entry:%s:", stream))
		fromPrefix := []byte(fmt.Sprintf("entry:%s:%020d:", stream, from.UnixMilli()))
		toPrefix := fmt.Sprintf("entry:%s:%020d:", stream, to.UnixMilli())

		opts := badger.DefaultIteratorOptions
		opts.Prefix = streamPrefix
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(fromPrefix); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())

			if key < toPrefix {
				timeKeys = append(timeKeys, key)
			} else {
				break
			}
		}

		var filteredKeys []string
		if filter != nil {
			fieldsIndex := newFieldsIndex(txn, stream)
			fieldKeys, err := fieldsIndex.Filter(filter)
			if err != nil {
				return &QueryError{Operation: "fields-index", Reason: err}
			}

			filteredKeys = intersectKeysList(timeKeys, fieldKeys)
		} else {
			filteredKeys = timeKeys
		}

		sort.Sort(sort.Reverse(sort.StringSlice(filteredKeys)))

		for _, key := range filteredKeys {
			item, err := txn.Get([]byte(key))
			if err != nil {
				return &QueryError{Operation: "query", Reason: err}
			}

			var entry LogEntry
			err = item.Value(func(val []byte) error {
				if err := json.Unmarshal(val, &entry); err != nil {
					return &UnmarshalError{Reason: err}
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

func (s *Storage) Purge(stream string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		prefixes := []string{
			fmt.Sprintf("entry:%s:", stream),
			fmt.Sprintf("index:%s:", stream),
		}

		for _, prefix := range prefixes {
			opts := badger.DefaultIteratorOptions
			opts.Prefix = []byte(prefix)
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				if err := txn.Delete(it.Item().Key()); err != nil {
					return &PersistError{Operation: "purge", Reason: err}
				}
			}
		}

		if err := txn.Delete([]byte(fmt.Sprintf("stream:%s", stream))); err != nil {
			return &PersistError{Operation: "purge", Reason: err}
		}

		return nil
	})
}

func (s *Storage) ListStreams() ([]string, error) {
	streams := []string{}

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("stream:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			stream := string(it.Item().Key()[7:])
			streams = append(streams, stream)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
