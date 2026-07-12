package transactions

import (
	"fmt"

	"time"

	"encoding/base64"
	"encoding/json"
	"strings"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

type fieldIndex struct {
	txn       kv.MutationTx
	keyPrefix kv.Key
}

// IndexField back-fills the inverted index for an already-populated stream: it
// scans every record, and for each one records its value of the given field,
// carrying over the entry's remaining retention so the index key expires with
// the entry.
func IndexField(txn kv.MutationTx, stream string, field string) error {
	for pair := range txn.IterPairs(kv.Key{"entry", stream}, kv.KeyRange{}) {
		key := pair.Key()

		var entry models.LogRecord
		if err := json.Unmarshal(pair.Value(), &entry); err != nil {
			return fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
		}

		ts := pair.ExpiresAt()
		retentionTime := int64(0)
		if ts != 0 {
			retentionTime = int64(ts) - time.Now().Unix()
		}

		index := newFieldIndex(txn, stream, field, entry.Fields[field])
		if err := index.AddKey(key, retentionTime); err != nil {
			return fmt.Errorf(
				"could not index field '%s' for entry '%s': %w",
				field, key, err,
			)
		}
	}

	return nil
}

// UnindexField drops the entire inverted index for a field by deleting every
// "index:<stream>:field:<field>:*" key.
func UnindexField(txn kv.MutationTx, stream string, field string) error {
	for key := range txn.IterKeys(kv.Key{"index", stream, "field", field}, kv.KeyRange{}) {
		if err := txn.Clear(key); err != nil {
			return fmt.Errorf("could not delete index key '%s': %w", key, err)
		}
	}

	return nil
}

// Distinct returns, per indexed field, the set of distinct values present in a
// stream. It reads them straight from the index keys (decoding the base64 value
// segment) without ever touching the log records themselves.
func Distinct(txn kv.QueryTx, stream string) (map[string][]string, error) {
	indices := make(map[string][]string)
	seenValuesPerField := make(map[string]map[string]struct{})

	for key := range txn.IterKeys(kv.Key{"index", stream, "field"}, kv.KeyRange{}) {
		if len(key) != 6 {
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
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))

	return &fieldIndex{
		txn:       txn,
		keyPrefix: kv.Key{"index", stream, "field", field, encodedValue},
	}
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
		return fmt.Errorf(
			"could not add field index key '%s' for entry '%s': %w",
			indexKey, entryKey, err,
		)
	}

	return nil
}

// purgeEntryFromFieldIndex removes every inverted-index reference to a deleted
// entry by scanning the stream's field index and dropping the keys whose
// trailing entry-key segment matches it.
func purgeEntryFromFieldIndex(txn kv.MutationTx, stream string, key kv.Key) error {
	for indexKey := range txn.IterKeys(kv.Key{"index", stream, "field"}, kv.KeyRange{}) {
		if indexKey.HasSuffix(key) {
			if err := txn.Clear(indexKey); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from field index: %w",
					strings.Join(key, ":"), err,
				)
			}
		}
	}

	return nil
}

// encodeIndexingMap base64-encodes the requested filter values so they line up
// with how values are encoded inside the index keys.
func encodeIndexingMap(indexing map[string][]string) map[string][]string {
	encodedMap := make(map[string][]string, len(indexing))

	for field, values := range indexing {
		encodedValues := make([]string, len(values))

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
