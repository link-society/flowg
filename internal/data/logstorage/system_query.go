package logstorage

import (
	"context"
	"fmt"
	"log/slog"

	"encoding/json"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type QuerySystem struct {
	storage *Storage
}

func NewQuerySystem(storage *Storage) *QuerySystem {
	return &QuerySystem{storage: storage}
}

func (sys *QuerySystem) FetchLogs(
	ctx context.Context,
	stream string,
	from, to time.Time,
	filter Filter,
) ([]LogEntry, error) {
	results := []LogEntry{}

	err := sys.storage.db.View(func(txn *badger.Txn) error {
		keys := fetchKeysByTime(txn, stream, from, to)
		keys = filterKeysByIndex(txn, stream, keys, filter)
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))

		for _, key := range keys {
			slog.DebugContext(
				ctx,
				"Fetching log entry",
				"channel", "storage",
				"stream", stream,
				"key", key,
			)

			entry, err := fetchEntry(txn, stream, key)
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

func filterKeysByIndex(
	txn *badger.Txn,
	stream string,
	allKeys []string,
	filter Filter,
) []string {
	allKeysMap := sliceToMap(allKeys)

	var evaluate func(f Filter) map[string]struct{}
	evaluate = func(f Filter) map[string]struct{} {
		switch f := f.(type) {
		case *AndFilter:
			keys := evaluate(f.Filters[0])

			for _, subFilter := range f.Filters[1:] {
				subKeys := evaluate(subFilter)
				keys = intersectKeysMap(keys, subKeys)
			}

			return keys

		case *OrFilter:
			keys := map[string]struct{}{}

			for _, subFilter := range f.Filters {
				subKeys := evaluate(subFilter)
				keys = unionKeysMap(keys, subKeys)
			}

			return keys

		case *NotFilter:
			keys := evaluate(f.Filter)
			return differenceKeysMap(allKeysMap, keys)

		case *FieldExact:
			fieldIndex := newFieldIndex(txn, stream, f.Field, f.Value)

			keys := map[string]struct{}{}

			fieldIndex.IterKeys(func(key string) {
				if _, found := allKeysMap[key]; found {
					keys[key] = struct{}{}
				}
			})

			return keys

		case *FieldIn:
			keys := map[string]struct{}{}

			for _, value := range f.Values {
				fieldIndex := newFieldIndex(txn, stream, f.Field, value)

				fieldIndex.IterKeys(func(key string) {
					if _, found := allKeysMap[key]; found {
						keys[key] = struct{}{}
					}
				})
			}

			return keys

		default:
			return map[string]struct{}{}
		}
	}

	keys := evaluate(filter)
	return mapToSlice(keys)
}

func fetchKeysByTime(txn *badger.Txn, stream string, from, to time.Time) []string {
	keys := []string{}

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
			keys = append(keys, key)
		} else {
			break
		}
	}

	return keys
}

func fetchEntry(txn *badger.Txn, stream string, key string) (LogEntry, error) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return LogEntry{}, fmt.Errorf(
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

	return entry, err
}
