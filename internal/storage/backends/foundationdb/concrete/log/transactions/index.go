package transactions

import (
	"fmt"

	"encoding/base64"
	"encoding/json"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/models"
)

type fieldIndex struct {
	txn       fdb.Transaction
	keyPrefix []byte
}

// IndexField back-fills the inverted index for an already-populated stream: it
// scans every record, and for each one records its value of the given field.
func IndexField(txn fdb.Transaction, stream, field string) error {
	streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))

	ri := txn.GetRange(streamEntrySS, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()
		key := append([]byte(nil), kv.Key...)

		var entry models.LogRecord
		if err := json.Unmarshal(kv.Value, &entry); err != nil {
			return fmt.Errorf("could not unmarshal log entry: %w", err)
		}

		index := newFieldIndex(txn, stream, field, entry.Fields[field])
		if err := index.AddKey(key); err != nil {
			return fmt.Errorf(
				"could not index field '%s' for entry: %w",
				field, err,
			)
		}
	}

	return nil
}

// UnindexField drops the entire inverted index for a field by clearing the
// <stream>:<field> sub-space in the index subspace.
func UnindexField(txn fdb.Transaction, stream, field string) error {
	streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))
	fieldIndexSS := streamIndexSS.Sub(subspace.FromBytes([]byte(field)))
	txn.ClearRange(fieldIndexSS)

	return nil
}

// Distinct returns, per indexed field, the set of distinct values present in a
// stream. It reads them straight from the index keys (decoding the base64 value
// segment) without ever touching the log records themselves.
func Distinct(txn fdb.ReadTransaction, stream string) (map[string][]string, error) {
	indices := make(map[string][]string)
	seenValuesPerField := make(map[string]map[string]struct{})

	streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))

	ri := txn.GetRange(streamIndexSS, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()

		field, encodedValue, ok := parseIndexKey(kv.Key)
		if !ok {
			continue
		}

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

// parseIndexKey extracts field and base64-encoded value from an index key.
//
// The key is structured as subspace hierarchy:
//
//	<indexSS prefix>
//	<stream subspace payload: 1-byte-len + stream bytes>
//	<field subspace payload:  1-byte-len + field bytes>
//	<base64 subspace payload: 1-byte-len + b64 bytes>
//	<tuple-encoded entry key>
//
// Each subspace.FromBytes([]byte(x)) appends len(x) as 1 byte then x bytes.
func parseIndexKey(key fdb.Key) (field string, encodedValue string, ok bool) {
	rest := key[len(indexSS.Bytes()):]
	if len(rest) < 3 {
		return "", "", false
	}

	// Skip stream: 1 byte length + stream bytes.
	streamLen := int(rest[0])
	rest = rest[1+streamLen:]
	if len(rest) < 3 {
		return "", "", false
	}

	// Read field: 1 byte length + field bytes.
	fieldLen := int(rest[0])
	if fieldLen < 1 || 1+fieldLen > len(rest) {
		return "", "", false
	}
	field = string(rest[1 : 1+fieldLen])
	rest = rest[1+fieldLen:]
	if len(rest) < 3 {
		return "", "", false
	}

	// Read base64 value: 1 byte length + b64 bytes.
	b64Len := int(rest[0])
	if b64Len < 1 || 1+b64Len > len(rest) {
		return "", "", false
	}
	encodedValue = string(rest[1 : 1+b64Len])

	return field, encodedValue, true
}

// newFieldIndex builds the index key prefix for one (stream, field, value)
// triple, base64-encoding the value as it appears in the stored keys.
func newFieldIndex(txn fdb.Transaction, stream, field, value string) *fieldIndex {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))

	// Build a subspace for this (stream, field, encodedValue) combination.
	streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))
	fieldIndexSS := streamIndexSS.Sub(subspace.FromBytes([]byte(field)))
	valueIndexSS := fieldIndexSS.Sub(subspace.FromBytes([]byte(encodedValue)))

	return &fieldIndex{
		txn:       txn,
		keyPrefix: valueIndexSS.Bytes(),
	}
}

// AddKey records that entryKey carries this (stream, field, value) by setting
// the index key: <valueSubspace><packed(entryKey)> = empty.
func (index *fieldIndex) AddKey(entryKey []byte) error {
	indexKey := subspace.FromBytes(index.keyPrefix).Pack(tuple.Tuple{entryKey})
	index.txn.Set(fdb.Key(indexKey), nil)
	return nil
}

// IterKeys yields the entry keys recorded under this (stream, field, value)
// index prefix.
func (index *fieldIndex) IterKeys(txn fdb.ReadTransaction, fn func(key []byte)) {
	valueSub := subspace.FromBytes(index.keyPrefix)

	ri := txn.GetRange(valueSub, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()

		tpl, err := tuple.Unpack(kv.Key[len(valueSub.Bytes()):])
		if err != nil || len(tpl) < 1 {
			continue
		}
		entryKey, ok := tpl[0].([]byte)
		if !ok {
			continue
		}
		fn(entryKey)
	}
}

// purgeEntryFromFieldIndex removes every inverted-index reference to a given
// entry by scanning the stream's index subspace and deleting keys whose
// last tuple element (the referenced entry key) matches.
func purgeEntryFromFieldIndex(txn fdb.Transaction, stream string, key []byte) error {
	streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))

	ri := txn.GetRange(streamIndexSS, fdb.RangeOptions{}).Iterator()
	for ri.Advance() {
		kv := ri.MustGet()

		tpl, err := tuple.Unpack(kv.Key[len(streamIndexSS.Bytes()):])
		if err != nil || len(tpl) < 1 {
			continue
		}
		refKey, ok := tpl[len(tpl)-1].([]byte)
		if !ok {
			continue
		}

		if string(refKey) == string(key) {
			txn.Clear(fdb.Key(kv.Key))
		}
	}

	return nil
}

// matchesIndexingForKey reports whether an entry satisfies the requested field
// filter: for every field the entry must match at least one of the requested
// values. Values are OR-ed within a field and the fields are AND-ed together.
func matchesIndexingForKey(
	txn fdb.ReadTransaction,
	stream string,
	entryKey string,
	encodedIndexing map[string][]string,
) (bool, error) {
	for field, encodedValues := range encodedIndexing {
		matched := false

		for _, encodedValue := range encodedValues {
			streamIndexSS := indexSS.Sub(subspace.FromBytes([]byte(stream)))
			fieldIndexSS := streamIndexSS.Sub(subspace.FromBytes([]byte(field)))
			valueIndexSS := fieldIndexSS.Sub(subspace.FromBytes([]byte(encodedValue)))
			indexKey := valueIndexSS.Pack(tuple.Tuple{[]byte(entryKey)})

			val, err := txn.Get(fdb.Key(indexKey)).Get()
			if err != nil {
				return false, fmt.Errorf(
					"could not check existence of index key: %w", err,
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
