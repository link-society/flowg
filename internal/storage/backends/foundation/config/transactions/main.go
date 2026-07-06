package transactions

import (
	"fmt"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

func buildKeyPrefix(keySpace subspace.Subspace, itemType string) subspace.Subspace {
	return keySpace.Sub(itemType)
}

func packKey(keySpace subspace.Subspace, itemType string, name string) fdb.Key {
	return buildKeyPrefix(keySpace, itemType).Pack(tuple.Tuple{name})
}

func unpackKey(keySpace subspace.Subspace, itemType string, key fdb.Key) (string, error) {
	keyTuple, err := buildKeyPrefix(keySpace, itemType).Unpack(key)
	if err != nil {
		return "", err
	}

	if len(keyTuple) < 1 {
		return "", nil
	}

	name, ok := keyTuple[0].(string)
	if !ok {
		return "", nil
	}

	return name, nil
}

// ListItems returns the names of every item of the given type by scanning the
// `<keySpace>/<itemType>/` subspace; only the keys carry the names, so values
// are skipped.
func ListItems(txn fdb.ReadTransaction, keySpace subspace.Subspace, itemType string) ([]string, error) {
	prefix := buildKeyPrefix(keySpace, itemType)
	it := txn.GetRange(prefix, fdb.RangeOptions{}).Iterator()

	var items []string
	for it.Advance() {
		kv, err := it.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to list items of type '%s': %w", itemType, err)
		}

		name, err := unpackKey(keySpace, itemType, kv.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to list items of type '%s': %w", itemType, err)
		}

		if name == "" {
			continue
		}

		items = append(items, name)
	}

	return items, nil
}

// ReadItem returns the raw, already-serialized value stored at
// `<keySpace>/<itemType>/<name>`, missing keys are returned as nil with no
// error.
func ReadItem(txn fdb.ReadTransaction, keySpace subspace.Subspace, itemType string, name string) ([]byte, error) {
	key := packKey(keySpace, itemType, name)
	v, err := txn.Get(key).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to read item '%s' of type '%s': %w", name, itemType, err)
	}
	return v, nil
}

// WriteItem stores content verbatim at `<keySpace>/<itemType>/<name>`, creating
// or overwriting the item.
func WriteItem(txn fdb.Transaction, keySpace subspace.Subspace, itemType string, name string, content []byte) error {
	key := packKey(keySpace, itemType, name)
	txn.Set(key, content)
	return nil
}

// DeleteItem removes the item at `<keySpace>/<itemType>/<name>`, if it exists.
// Deleting a key that does not exist is a no-op.
func DeleteItem(txn fdb.Transaction, keySpace subspace.Subspace, itemType string, name string) error {
	key := packKey(keySpace, itemType, name)
	txn.Clear(key)
	return nil
}
