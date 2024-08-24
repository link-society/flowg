package logstorage

import (
	"encoding/base64"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type fieldIndex struct {
	txn       *badger.Txn
	keyPrefix []byte
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

func (index *fieldIndex) AddKey(entryKey []byte) error {
	indexKey := []byte(fmt.Sprintf("%s%s", index.keyPrefix, entryKey))
	return index.txn.Set(indexKey, []byte{})
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
