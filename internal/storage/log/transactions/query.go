package transactions

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/ffi/filterdsl"
)

func FetchLogs(
	txn *badger.Txn,
	stream string,
	from, to time.Time,
	filter filterdsl.Filter,
) ([]models.LogRecord, error) {
	results := []models.LogRecord{}

	streamConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return nil, err
	}

	keys := fetchKeysByTime(txn, stream, from, to)
	if filter != nil {
		keys = filterKeysByIndex(txn, stream, streamConfig, keys, filter)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, key := range keys {
		entry, err := fetchRecord(txn, stream, key)
		if err != nil {
			return nil, err
		}

		if filter == nil || filter.Evaluate(&entry) {
			results = append(results, entry)
		}
	}

	return results, nil
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

func filterKeysByIndex(
	txn *badger.Txn,
	stream string,
	streamConfig models.StreamConfig,
	allKeys []string,
	filter filterdsl.Filter,
) []string {
	allKeysMap := sliceToMap(allKeys)

	var evaluate func(f filterdsl.Filter) map[string]struct{}
	evaluate = func(f filterdsl.Filter) map[string]struct{} {
		switch f := f.(type) {
		case *filterdsl.FilterAnd:
			keys := evaluate(f.Filters[0])

			for _, subFilter := range f.Filters[1:] {
				subKeys := evaluate(subFilter)
				keys = intersectKeysMap(keys, subKeys)
			}

			return keys

		case *filterdsl.FilterOr:
			keys := map[string]struct{}{}

			for _, subFilter := range f.Filters {
				subKeys := evaluate(subFilter)
				keys = unionKeysMap(keys, subKeys)
			}

			return keys

		case *filterdsl.FilterNot:
			switch sub := f.Filter.(type) {
			case *filterdsl.FilterMatchField:
				if streamConfig.IsFieldIndexed(sub.Field) {
					keys := evaluate(f.Filter)
					return differenceKeysMap(allKeysMap, keys)
				} else {
					return allKeysMap
				}

			case *filterdsl.FilterMatchFieldList:
				if streamConfig.IsFieldIndexed(sub.Field) {
					keys := evaluate(f.Filter)
					return differenceKeysMap(allKeysMap, keys)
				} else {
					return allKeysMap
				}

			default:
				keys := evaluate(f.Filter)
				return differenceKeysMap(allKeysMap, keys)
			}

		case *filterdsl.FilterMatchField:
			if streamConfig.IsFieldIndexed(f.Field) {
				fieldIndex := newFieldIndex(txn, stream, f.Field, f.Value)

				keys := map[string]struct{}{}

				fieldIndex.IterKeys(func(key string) {
					if _, found := allKeysMap[key]; found {
						keys[key] = struct{}{}
					}
				})

				return keys
			} else {
				return allKeysMap
			}

		case *filterdsl.FilterMatchFieldList:
			if streamConfig.IsFieldIndexed(f.Field) {
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
			} else {
				return allKeysMap
			}

		default:
			return map[string]struct{}{}
		}
	}

	keys := evaluate(filter)
	return mapToSlice(keys)
}

func fetchRecord(txn *badger.Txn, stream string, key string) (models.LogRecord, error) {
	item, err := txn.Get([]byte(key))
	if err != nil {
		return models.LogRecord{}, fmt.Errorf(
			"could not fetch log entry '%s' from stream '%s': %w",
			key, stream, err,
		)
	}

	var record models.LogRecord
	err = item.Value(func(val []byte) error {
		if err := json.Unmarshal(val, &record); err != nil {
			return fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
		}

		return nil
	})

	return record, err
}

func intersectKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		if _, exists := b[key]; exists {
			result[key] = struct{}{}
		}
	}

	return result
}

func unionKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		result[key] = struct{}{}
	}

	for key := range b {
		result[key] = struct{}{}
	}

	return result
}

func differenceKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		if _, exists := b[key]; !exists {
			result[key] = struct{}{}
		}
	}

	return result
}

func sliceToMap(slice []string) map[string]struct{} {
	result := map[string]struct{}{}

	for _, key := range slice {
		result[key] = struct{}{}
	}

	return result
}

func mapToSlice(m map[string]struct{}) []string {
	result := []string{}

	for key := range m {
		result = append(result, key)
	}

	return result
}
