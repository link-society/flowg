package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

func ListItems(txn *badger.Txn, itemType string) ([]string, error) {
	items := []string{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("%s:", itemType))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		items = append(items, string(key[len(itemType)+1:]))
	}

	return items, nil
}

func ReadItem(txn *badger.Txn, itemType string, name string) ([]byte, error) {
	item, err := txn.Get([]byte(fmt.Sprintf("%s:%s", itemType, name)))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get %s '%s': %w", itemType, name, err)
	}

	content, err := item.ValueCopy(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s '%s': %w", itemType, name, err)
	}

	return content, nil
}

func WriteItem(txn *badger.Txn, itemType string, name string, content []byte) error {
	key := []byte(fmt.Sprintf("%s:%s", itemType, name))
	err := txn.Set(key, content)
	if err != nil {
		return fmt.Errorf("failed to write %s '%s': %w", itemType, name, err)
	}

	return nil
}

func DeleteItem(txn *badger.Txn, itemType string, name string) error {
	key := []byte(fmt.Sprintf("%s:%s", itemType, name))
	err := txn.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete %s '%s': %w", itemType, name, err)
	}

	return nil
}
