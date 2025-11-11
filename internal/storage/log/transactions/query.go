package transactions

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/langs/filtering"
)

func FetchLogs(
	txn *badger.Txn,
	stream string,
	from, to time.Time,
	filter filtering.Filter,
) ([]models.LogRecord, error) {
	results := []models.LogRecord{}

	keys := fetchKeysByTime(txn, stream, from, to)
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	for _, key := range keys {
		entry, err := fetchRecord(txn, stream, key)
		if err != nil {
			return nil, err
		}

		if filter == nil {
			results = append(results, entry)
		} else {
			matches, err := filter.Evaluate(&entry)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to evaluate filter for log entry '%s': %w",
					key, err,
				)
			}

			if matches {
				results = append(results, entry)
			}
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
