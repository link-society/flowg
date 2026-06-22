package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func ListItems(txn *badger.Txn, itemType string) ([]string, error) {
	items := []string{}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = fmt.Appendf(nil, "%s:", itemType)
	it := txn.NewIterator(opts)
	defer it.Close()

	prefixLen := len(itemType) + 1

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		name := string(item.Key()[prefixLen:])

		err := item.Value(func(val []byte) error {
			env, err := lww.Unmarshal(val)
			if err != nil {
				return err
			}
			if !env.Deleted {
				items = append(items, name)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to read %s '%s': %w", itemType, name, err)
		}
	}

	return items, nil
}

func ReadItem(txn *badger.Txn, itemType string, name string) ([]byte, error) {
	key := fmt.Appendf(nil, "%s:%s", itemType, name)

	env, found, err := lww.Read(txn, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get %s '%s': %w", itemType, name, err)
	}
	if !found {
		return nil, fmt.Errorf("%s '%s': %w", itemType, name, badger.ErrKeyNotFound)
	}

	return env.Payload, nil
}

func WriteItem(txn *badger.Txn, itemType string, name string, content []byte, ts hlc.Timestamp) (bool, error) {
	key := fmt.Appendf(nil, "%s:%s", itemType, name)

	applied, err := lww.Apply(txn, key, lww.Envelope{Timestamp: ts, Payload: content})
	if err != nil {
		return false, fmt.Errorf("failed to write %s '%s': %w", itemType, name, err)
	}

	return applied, nil
}

func DeleteItem(txn *badger.Txn, itemType string, name string, ts hlc.Timestamp) (bool, error) {
	key := fmt.Appendf(nil, "%s:%s", itemType, name)

	applied, err := lww.Apply(txn, key, lww.Envelope{Timestamp: ts, Deleted: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete %s '%s': %w", itemType, name, err)
	}

	return applied, nil
}
