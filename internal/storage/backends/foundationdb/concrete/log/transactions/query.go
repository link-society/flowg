package transactions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/langs/filtering"
)

// FetchLogs returns the records of a stream between two timestamps, newest
// first. It gathers the candidate keys in the time window, narrows them to those
// present in the requested field indexes, then decodes each surviving record and
// keeps the ones the filter accepts.
func FetchLogs(
	txn fdb.ReadTransaction,
	stream string,
	from, to time.Time,
	filter filtering.Filter,
	indexing map[string][]string,
) ([]models.LogRecord, error) {
	results := []models.LogRecord{}

	keys, err := fetchKeysByTime(txn, stream, from, to)
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	var keysForIndex []string
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
					"failed to evaluate filter for log entry: %w", err,
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
// [from, to). Because the timestamp is packed as a tuple element that sorts
// numerically, it can seek straight to the lower bound and stop as soon as a
// key sorts past the upper bound instead of scanning the whole stream.
func fetchKeysByTime(txn fdb.ReadTransaction, stream string, from, to time.Time) ([]string, error) {
	keys := []string{}

	streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))
	fromPadded := fmt.Sprintf("%020d", from.UnixMilli())
	toPadded := fmt.Sprintf("%020d", to.UnixMilli())

	startKey := fdb.Key(streamEntrySS.Pack(tuple.Tuple{fromPadded}))
	endKey := fdb.Key(streamEntrySS.Pack(tuple.Tuple{toPadded}))

	kr := fdb.KeyRange{Begin: startKey, End: endKey}
	ri := txn.GetRange(kr, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()
		keys = append(keys, string(kv.Key))
	}

	return keys, nil
}

// fetchRecord loads and JSON-decodes a single log record by its key.
// encodeIndexingMap base64-encodes the requested filter values so they line up
// with how values are encoded inside the index keys.
func encodeIndexingMap(indexing map[string][]string) map[string][]string {
	encodedMap := make(map[string][]string, len(indexing))

	for field, values := range indexing {
		encodedValues := make([]string, 0, len(values))

		for _, value := range values {
			encodedValues = append(
				encodedValues,
				base64.StdEncoding.EncodeToString([]byte(value)),
			)
		}

		encodedMap[field] = encodedValues
	}

	return encodedMap
}

func fetchRecord(txn fdb.ReadTransaction, stream string, key string) (models.LogRecord, error) {
	val, err := txn.Get(fdb.Key(key)).Get()
	if err != nil {
		return models.LogRecord{}, fmt.Errorf(
			"could not fetch log entry from stream '%s': %w",
			stream, err,
		)
	}

	if val == nil {
		return models.LogRecord{}, fmt.Errorf(
			"log entry not found in stream '%s'", stream,
		)
	}

	var record models.LogRecord
	if err := json.Unmarshal(val, &record); err != nil {
		return models.LogRecord{}, fmt.Errorf(
			"could not unmarshal log entry from stream '%s': %w",
			stream, err,
		)
	}

	return record, nil
}
