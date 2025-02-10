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

func newFieldIndex(txn *badger.Txn, stream, field, value string) *fieldIndex {
	encodedValue := base64.StdEncoding.EncodeToString([]byte(value))

	return &fieldIndex{
		txn: txn,
		keyPrefix: []byte(fmt.Sprintf(
			"index:%s:field:%s:%s:",
			stream, field, encodedValue,
		)),
	}
}

func (index *fieldIndex) AddKey(entryKey []byte, retentionTime int64) error {
	indexKey := []byte(fmt.Sprintf("%s%s", index.keyPrefix, entryKey))
	entry := badger.NewEntry(indexKey, []byte{})

	if retentionTime > 0 {
		entry = entry.WithTTL(time.Duration(retentionTime) * time.Second)
	}

	return index.txn.SetEntry(entry)
}

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
