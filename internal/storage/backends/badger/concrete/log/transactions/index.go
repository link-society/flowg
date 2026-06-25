package transactions

import (
	"fmt"
	"strings"

	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
)

type fieldIndex struct {
	txn       *badger.Txn
	keyPrefix []byte
}

// IndexField back-fills the inverted index for an already-populated stream: it
// scans every record, and for each one records its value of the given field,
// carrying over the entry's remaining retention so the index key expires with
// the entry.
func IndexField(txn *badger.Txn, stream, field string) error {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte("entry:" + stream + ":")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		key := item.KeyCopy(nil)

		var entry models.LogRecord
		err := item.Value(func(val []byte) error {
			if err := json.Unmarshal(val, &entry); err != nil {
				return fmt.Errorf("could not unmarshal log entry '%s': %w", key, err)
			}

			return nil
		})
		if err != nil {
			return err
		}

		ts := item.ExpiresAt()
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
func UnindexField(txn *badger.Txn, stream, field string) error {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("index:%s:field:%s:", stream, field))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		indexKey := it.Item().KeyCopy(nil)

		if err := txn.Delete(indexKey); err != nil {
			return fmt.Errorf("could not delete index key '%s': %w", indexKey, err)
		}
	}

	return nil
}

// Distinct returns, per indexed field, the set of distinct values present in a
// stream. It reads them straight from the index keys (decoding the base64 value
// segment) without ever touching the log records themselves.
func Distinct(txn *badger.Txn, stream string) (map[string][]string, error) {
	indices := make(map[string][]string)
	seenValuesPerField := make(map[string]map[string]struct{})

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = fmt.Appendf(nil, "index:%s:field:", stream)
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		indexKey := it.Item().Key()
		parts := strings.SplitN(string(indexKey), ":", 6)
		if len(parts) != 6 {
			continue
		}

		field := parts[3]
		encodedValue := parts[4]

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
func newFieldIndex(txn *badger.Txn, stream, field, value string) *fieldIndex {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))

	return &fieldIndex{
		txn: txn,
		keyPrefix: fmt.Appendf(nil,
			"index:%s:field:%s:%s:",
			stream, field, encodedValue,
		),
	}
}

// AddKey records that entryKey carries this (stream, field, value), giving the
// index key the entry's remaining TTL so the two expire together.
func (index *fieldIndex) AddKey(entryKey []byte, retentionTime int64) error {
	indexKey := []byte(fmt.Sprintf("%s%s", index.keyPrefix, entryKey))
	entry := badger.NewEntry(indexKey, []byte{})

	if retentionTime > 0 {
		entry = entry.WithTTL(time.Duration(retentionTime) * time.Second)
	}

	return index.txn.SetEntry(entry)
}

// IterKeys yields the entry keys recorded under this (stream, field, value)
// index prefix.
func (index *fieldIndex) IterKeys(fn func(key string)) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = index.keyPrefix
	it := index.txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		indexKey := it.Item().Key()
		entryKey := string(indexKey[len(index.keyPrefix):])
		fn(entryKey)
	}
}

// purgeEntryFromFieldIndex removes every inverted-index reference to a deleted
// entry by scanning the stream's field index and dropping the keys whose
// trailing entry-key segment matches it.
func purgeEntryFromFieldIndex(txn *badger.Txn, stream string, key []byte) error {
	suffix := fmt.Sprintf(":%s", string(key))

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("index:%s:field:", stream))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		indexKey := string(item.Key())

		if strings.HasSuffix(indexKey, suffix) {
			if err := txn.Delete(item.KeyCopy(nil)); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from field index: %w",
					string(key), err,
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
	txn *badger.Txn,
	stream string,
	entryKey string,
	encodedIndexing map[string][]string,
) (bool, error) {
	for field, encodedValues := range encodedIndexing {
		matched := false

		for _, encodedValue := range encodedValues {
			indexKey := fmt.Sprintf(
				"index:%s:field:%s:%s:%s",
				stream, field, encodedValue, entryKey,
			)

			_, err := txn.Get([]byte(indexKey))
			if err == nil {
				matched = true
				break
			}
			if err != badger.ErrKeyNotFound {
				return false, fmt.Errorf(
					"could not check existence of index key '%s': %w",
					indexKey, err,
				)
			}
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}
