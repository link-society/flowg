package storage

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type fieldIndex struct {
	txn *badger.Txn
	key []byte
}

func newFieldIndex(txn *badger.Txn, stream, field, value string) *fieldIndex {
	return &fieldIndex{
		txn: txn,
		key: []byte(fmt.Sprintf("index:%s:field:%s:%s", stream, field, value)),
	}
}

func (index *fieldIndex) AddKey(entryKey []byte) error {
	fieldKeys, err := index.load()
	if err != nil {
		return err
	}

	fieldKeys = append(fieldKeys, string(entryKey))

	return index.store(fieldKeys)
}

func (index *fieldIndex) GetKeys() ([]string, error) {
	return index.load()
}

func (index *fieldIndex) load() ([]string, error) {
	fieldKeys := []string{}

	fieldIndexValue, err := index.txn.Get(index.key)
	if err == badger.ErrKeyNotFound {
		return fieldKeys, nil
	} else if err != nil {
		return nil, err
	}

	err = fieldIndexValue.Value(func(val []byte) error {
		return json.Unmarshal(val, &fieldKeys)
	})
	if err != nil {
		return nil, &UnmarshalError{Reason: err}
	}

	return fieldKeys, nil
}

func (index *fieldIndex) store(fieldKeys []string) error {
	fieldKeysEncoded, err := json.Marshal(fieldKeys)
	if err != nil {
		return &MarshalError{Reason: err}
	}

	if err := index.txn.Set(index.key, fieldKeysEncoded); err != nil {
		return err
	}

	return nil
}
