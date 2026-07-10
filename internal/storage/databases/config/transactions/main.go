package transactions

import (
	"fmt"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// ListItems returns the names of every item of the given type by scanning the
// "<itemType>" prefix; only the keys carry the names, so values are skipped.
func ListItems(txn kv.QueryTx, itemType string) ([]string, error) {
	items := []string{}

	for key := range txn.IterKeys(kv.Key{itemType}, kv.KeyRange{}) {
		itemName := key[len(key)-1]
		items = append(items, itemName)
	}

	return items, nil
}

// ReadItem returns the raw, already-serialized value stored at
// "<itemType>:<name>", wrapping a missing key as a typed not-found error.
func ReadItem(txn kv.QueryTx, itemType string, name string) ([]byte, error) {
	content, err := txn.Get(kv.Key{itemType, name})
	if err != nil {
		return nil, fmt.Errorf("failed to read %s %q: %w", itemType, name, err)
	}

	return content, nil
}

// WriteItem stores content verbatim at "<itemType>:<name>", creating or
// overwriting the item.
func WriteItem(txn kv.MutationTx, itemType string, name string, content []byte) error {
	err := txn.Set(kv.Key{itemType, name}, content)
	if err != nil {
		return fmt.Errorf("failed to write %s %q: %w", itemType, name, err)
	}

	return nil
}

// DeleteItem removes the "<itemType>:<name>" key; deleting a key that does not
// exist is a no-op.
func DeleteItem(txn kv.MutationTx, itemType string, name string) error {
	err := txn.Clear(kv.Key{itemType, name})
	if err != nil {
		return fmt.Errorf("failed to delete %s %q: %w", itemType, name, err)
	}

	return nil
}
