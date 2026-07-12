package transactions

import (
	"errors"
	"fmt"
	"log/slog"

	"time"

	"encoding/base64"
	"encoding/json"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

type fieldIndex struct {
	txn       kv.MutationTx
	stream    string
	field     string
	keyPrefix kv.Key
}

// IndexFieldBatch back-fills the inverted index for at most limit entries of an
// already-populated stream, resuming strictly after cursor (nil to start from
// the beginning). For each entry it records its value of the given field,
// carrying over the entry's remaining retention so the index key expires with
// the entry.
//
// It returns the last entry key it processed — to be passed as cursor on the
// next call — and how many entries it processed. A processed count below limit
// means the stream has been fully back-filled. Splitting the work this way keeps
// each transaction within the backend's size and time limits.
func IndexFieldBatch(
	txn kv.MutationTx,
	stream string,
	field string,
	cursor kv.Key,
	limit int,
) (kv.Key, int, error) {
	var last kv.Key
	processed := 0

	for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{From: cursor}) {
		key := pair.Key()

		// KeyRange.From is inclusive on some backends and exclusive on others;
		// skip the cursor itself so a batch always makes forward progress.
		if cursor != nil && keysEqual(key, cursor) {
			continue
		}

		if processed >= limit {
			break
		}

		var entry models.LogRecord
		if err := json.Unmarshal(pair.Value(), &entry); err != nil {
			return last, processed, fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
		}

		ts := pair.ExpiresAt()
		retentionTime := int64(0)
		if ts != 0 {
			retentionTime = int64(ts) - time.Now().Unix()
		}

		index := newFieldIndex(txn, stream, field, entry.Fields[field])
		if err := index.AddKey(key, retentionTime); err != nil {
			return last, processed, fmt.Errorf(
				"could not index field '%s' for entry '%s': %w",
				field, key, err,
			)
		}

		last = key
		processed++
	}

	return last, processed, nil
}

// UnindexFieldBatch drops up to limit keys of a field's inverted index
// ("index:<stream>:field:<field>:*") and reports how many it deleted. A count
// below limit means the index has been fully removed. Callers loop until it
// returns zero so the work stays within the backend's transaction limits.
func UnindexFieldBatch(txn kv.MutationTx, stream string, field string, limit int) (int, error) {
	var keys []kv.Key
	for key := range txn.IterKeys(kv.Key{"index", stream, "field", field}, kv.KeyRange{}) {
		keys = append(keys, key)
		if len(keys) >= limit {
			break
		}
	}

	for _, key := range keys {
		if err := txn.Clear(key); err != nil {
			return len(keys), fmt.Errorf("could not delete index key '%s': %w", key, err)
		}
	}

	return len(keys), nil
}

// Distinct returns, per indexed field, the set of distinct values present in a
// stream. It reads them straight from the index keys (decoding the base64 value
// segment) without ever touching the log records themselves.
func Distinct(txn kv.QueryTx, stream string) (map[string][]string, error) {
	indices := make(map[string][]string)
	seenValuesPerField := make(map[string]map[string]struct{})

	for key := range txn.IterKeys(kv.Key{"index", stream, "field"}, kv.KeyRange{}) {
		// index keys are "index:<stream>:field:<field>:<value>" followed by the
		// full entry key, so anything shorter than that 5-segment prefix is not a
		// value index key we can decode.
		if len(key) < 5 {
			continue
		}

		field := key[3]
		encodedValue := key[4]

		if _, exists := seenValuesPerField[field]; !exists {
			seenValuesPerField[field] = make(map[string]struct{})
		}

		if _, seen := seenValuesPerField[field][encodedValue]; !seen {
			decodedValueBytes, err := base64.StdEncoding.DecodeString(encodedValue)
			if err != nil {
				return nil, fmt.Errorf(
					"could not decode base64 value '%s' for field '%s': %w",
					encodedValue, field, err,
				)
			}
			decodedValue := string(decodedValueBytes)

			indices[field] = append(indices[field], decodedValue)
			seenValuesPerField[field][encodedValue] = struct{}{}
		}
	}

	return indices, nil
}

// newFieldIndex builds the index key prefix for one (stream, field, value)
// triple, base64-encoding the value as it appears in the stored keys.
func newFieldIndex(txn kv.MutationTx, stream string, field string, value string) *fieldIndex {
	return &fieldIndex{
		txn:       txn,
		stream:    stream,
		field:     field,
		keyPrefix: fieldIndexPrefix(stream, field, value),
	}
}

// fieldIndexPrefix builds the "index:<stream>:field:<field>:<base64(value)>"
// prefix shared by every entry carrying that (stream, field, value).
func fieldIndexPrefix(stream string, field string, value string) kv.Key {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))
	return kv.Key{"index", stream, "field", field, encodedValue}
}

// keysEqual reports whether two composite keys have identical segments.
func keysEqual(a, b kv.Key) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}


// AddKey records that entryKey carries this (stream, field, value), giving the
// index key the entry's remaining TTL so the two expire together.
func (index *fieldIndex) AddKey(entryKey kv.Key, retentionTime int64) error {
	indexKey := append(index.keyPrefix, entryKey...)

	var err error
	if retentionTime > 0 {
		ttl := time.Duration(retentionTime) * time.Second
		err = index.txn.SetWithTTL(indexKey, []byte{}, ttl)
	} else {
		err = index.txn.Set(indexKey, []byte{})
	}
	if err != nil {
		// An indexed field value is embedded in the index key. A large value
		// (e.g. a message body) can push the key past the backend's size limit.
		// Skip indexing that value rather than rejecting the whole record: the
		// entry is still stored and queryable by time, it just won't be found by
		// an exact-match filter on this oversized value.
		if errors.Is(err, kv.ErrKeyTooLarge) {
			slog.Warn(
				"skipping oversized field index value",
				slog.String("channel", "logstorage"),
				slog.String("stream", index.stream),
				slog.String("field", index.field),
			)
			return nil
		}

		return fmt.Errorf(
			"could not add field index key '%s' for entry '%s': %w",
			indexKey, entryKey, err,
		)
	}

	return nil
}

// purgeEntryFromFieldIndex removes the inverted-index references of a deleted
// entry. Because a log record is immutable, its stored fields fully determine
// the index keys that can point at it, so each key is reconstructed and cleared
// directly (O(fields)) instead of scanning the whole field index (O(index
// size)). Clearing a key that was never indexed is a harmless no-op.
func purgeEntryFromFieldIndex(txn kv.MutationTx, stream string, entryKey kv.Key, record *models.LogRecord) error {
	for field, value := range record.Fields {
		prefix := fieldIndexPrefix(stream, field, value)
		indexKey := make(kv.Key, 0, len(prefix)+len(entryKey))
		indexKey = append(indexKey, prefix...)
		indexKey = append(indexKey, entryKey...)

		if err := txn.Clear(indexKey); err != nil {
			// A value too large to fit in an index key was never indexed (see
			// AddKey), so there is nothing to purge for it. Fields that are not
			// indexed at all simply have no key to clear (a harmless no-op).
			if errors.Is(err, kv.ErrKeyTooLarge) {
				continue
			}

			return fmt.Errorf(
				"could not delete key '%s' from field index of stream '%s': %w",
				indexKey, stream, err,
			)
		}
	}

	return nil
}

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

// matchesIndexingForKey reports whether an entry satisfies the requested field
// filter: for every field the entry must match at least one of the requested
// values (i.e. the corresponding index key must exist). Values are OR-ed within
// a field and the fields are AND-ed together.
func matchesIndexingForKey(
	txn kv.QueryTx,
	stream string,
	entryKey kv.Key,
	encodedIndexing map[string][]string,
) (bool, error) {
	for field, encodedValues := range encodedIndexing {
		matched := false

		for _, encodedValue := range encodedValues {
			indexKey := kv.Key{"index", stream, "field", field, encodedValue}
			indexKey = append(indexKey, entryKey...)

			val, err := txn.Get(indexKey)
			if err != nil {
				// The value was too large to index (see AddKey), so no entry can
				// carry it: treat it as a non-match rather than an error.
				if errors.Is(err, kv.ErrKeyTooLarge) {
					continue
				}

				return false, fmt.Errorf(
					"could not check existence of index key '%s': %w",
					indexKey, err,
				)
			}
			if val != nil {
				matched = true
				break
			}
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}
