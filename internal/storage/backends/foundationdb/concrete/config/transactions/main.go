package transactions

import (
	"errors"
	"fmt"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

var root subspace.Subspace

// Init initializes the package-level root subspace.
func Init(rootSub subspace.Subspace) {
	root = rootSub
}

// ListItems returns the names of every item of the given type by scanning the
// <root>/<itemType>/ prefix; only the keys carry the names, so values are
// skipped.
func ListItems(tr fdb.ReadTransaction, itemType string) ([]string, error) {
	items := []string{}
	sub := root.Sub(itemType)

	iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := sub.Unpack(kv.Key)
		if err != nil {
			continue
		}

		items = append(items, t[0].(string))
	}

	return items, nil
}

// ReadItem returns the raw, already-serialized value stored at
// <root>/<itemType>/<name>, wrapping a missing key as a typed not-found error.
func ReadItem(tr fdb.ReadTransaction, itemType string, name string) ([]byte, error) {
	val, err := tr.Get(root.Sub(itemType).Pack(tuple.Tuple{name})).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s '%s': %w", itemType, name, err)
	}
	if val == nil {
		return nil, fmt.Errorf("%s '%s': %w", itemType, name, errKeyNotFound)
	}

	return val, nil
}

// WriteItem stores content verbatim at <root>/<itemType>/<name>, creating or
// overwriting the item.
func WriteItem(tr fdb.Transaction, itemType string, name string, content []byte) error {
	tr.Set(root.Sub(itemType).Pack(tuple.Tuple{name}), content)
	return nil
}

// DeleteItem removes the <root>/<itemType>/<name> key; deleting a key that
// does not exist is a no-op.
func DeleteItem(tr fdb.Transaction, itemType string, name string) error {
	tr.Clear(root.Sub(itemType).Pack(tuple.Tuple{name}))
	return nil
}

// KeyNotFound reports whether an error returned by ReadItem was caused by a
// missing key, so callers can distinguish "item doesn't exist" from a genuine
// I/O error.
func KeyNotFound(err error) bool {
	if err == nil {
		return false
	}

	// FDB does not signal missing keys as a sentinel error — the value is
	// simply nil. Instead, we detect the wrapped message produced by ReadItem.
	return errors.Is(err, errKeyNotFound)
}

var errKeyNotFound = errors.New("key not found")
