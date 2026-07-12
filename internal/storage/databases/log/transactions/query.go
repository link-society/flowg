package transactions

import (
	"fmt"

	"sort"
	"time"

	"encoding/json"
	"strings"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"

	"link-society.com/flowg/internal/utils/langs/filtering"
)

// FetchLogs returns the records of a stream between two timestamps, newest
// first. It gathers the candidate keys in the time window, narrows them to those
// present in the requested field indexes, then decodes each surviving record and
// keeps the ones the filter accepts.
func FetchLogs(
	txn kv.QueryTx,
	stream string,
	from, to time.Time,
	filter filtering.Filter,
	indexing map[string][]string,
) ([]models.LogRecord, error) {
	results := []models.LogRecord{}

	keys := fetchKeysByTime(txn, stream, from, to)
	sort.Sort(sort.Reverse(kv.KeySlice(keys)))

	var keysForIndex []kv.Key
	encodedIndexing := encodeIndexingMap(indexing)

	for _, key := range keys {
		matched, err := matchesIndexingForKey(txn, stream, key, encodedIndexing)
		if err != nil {
			return nil, err
		}

		if matched {
			keysForIndex = append(keysForIndex, key)
		}
	}

	for _, key := range keysForIndex {
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
					strings.Join(key, ":"), err,
				)
			}

			if matches {
				results = append(results, entry)
			}
		}
	}

	return results, nil
}

// fetchKeysByTime collects the entry keys whose embedded timestamp falls in
// [from, to]. Because the timestamp is zero-padded in the key, it can seek
// straight to the lower bound and stop as soon as a key sorts past the upper
// bound instead of scanning the whole stream.
func fetchKeysByTime(txn kv.QueryTx, stream string, from, to time.Time) []kv.Key {
	keys := []kv.Key{}

	streamPrefix := kv.Key{"entry", stream}
	fromPrefix := kv.Key{"entry", stream, fmt.Sprintf("%020d", from.UnixMilli())}
	toPrefix := kv.Key{"entry", stream, fmt.Sprintf("%020d", to.UnixMilli())}

	for key := range txn.IterKeys(streamPrefix, kv.KeyRange{From: fromPrefix, To: toPrefix}) {
		keys = append(keys, key)
	}

	return keys
}

// fetchRecord loads and JSON-decodes a single log record by its key.
func fetchRecord(txn kv.QueryTx, stream string, key kv.Key) (models.LogRecord, error) {
	val, err := txn.Get(key)
	if err != nil {
		return models.LogRecord{}, fmt.Errorf(
			"could not fetch log entry '%s' from stream '%s': %w",
			strings.Join(key, ":"), stream, err,
		)
	}

	var record models.LogRecord
	if err := json.Unmarshal(val, &record); err != nil {
		return models.LogRecord{}, fmt.Errorf("could not unmarshal log entry '%s': %w", strings.Join(key, ":"), err)
	}

	return record, nil
}
