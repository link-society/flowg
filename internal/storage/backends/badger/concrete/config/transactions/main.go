package transactions

import (
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

// ListItems returns the names of every item of the given type by scanning the
// "<itemType>:" prefix; only the keys carry the names, so values are skipped.
func ListItems(txn *badger.Txn, itemType string) ([]string, error) {
	items := []string{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = fmt.Appendf(nil, "%s:", itemType)
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		items = append(items, string(key[len(itemType)+1:]))
	}

	return items, nil
}

// ReadItem returns the raw, already-serialized value stored at
// "<itemType>:<name>", wrapping a missing key as a typed not-found error.
func ReadItem(txn *badger.Txn, itemType string, name string) ([]byte, error) {
	item, err := txn.Get(fmt.Appendf(nil, "%s:%s", itemType, name))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("%s '%s': %w", itemType, name, err)
		}

		return nil, fmt.Errorf("failed to get %s '%s': %w", itemType, name, err)
	}

	content, err := item.ValueCopy(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s '%s': %w", itemType, name, err)
	}

	return content, nil
}

// WriteItem stores content verbatim at "<itemType>:<name>", creating or
// overwriting the item.
func WriteItem(txn *badger.Txn, itemType string, name string, content []byte) error {
	key := fmt.Appendf(nil, "%s:%s", itemType, name)
	err := txn.Set(key, content)
	if err != nil {
		return fmt.Errorf("failed to write %s '%s': %w", itemType, name, err)
	}

	return nil
}

// DeleteItem removes the "<itemType>:<name>" key; deleting a key that does not
// exist is a no-op.
func DeleteItem(txn *badger.Txn, itemType string, name string) error {
	key := fmt.Appendf(nil, "%s:%s", itemType, name)
	err := txn.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete %s '%s': %w", itemType, name, err)
	}

	return nil
}
