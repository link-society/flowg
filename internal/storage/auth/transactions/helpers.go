package transactions

import (
	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func setItem(txn *badger.Txn, key []byte, payload []byte, ts hlc.Timestamp) error {
	_, err := lww.Apply(txn, key, lww.Envelope{Timestamp: ts, Payload: payload})
	return err
}

func deleteItem(txn *badger.Txn, key []byte, ts hlc.Timestamp) error {
	_, err := lww.Apply(txn, key, lww.Envelope{Timestamp: ts, Deleted: true})
	return err
}

func getItem(txn *badger.Txn, key []byte) ([]byte, bool, error) {
	env, found, err := lww.Read(txn, key)
	if err != nil {
		return nil, false, err
	}
	return env.Payload, found, nil
}

func liveValue(item *badger.Item) ([]byte, bool, error) {
	var (
		payload []byte
		live    bool
	)

	err := item.Value(func(val []byte) error {
		env, err := lww.Unmarshal(val)
		if err != nil {
			return err
		}
		if env.Deleted {
			return nil
		}
		live = true
		payload = append([]byte(nil), env.Payload...)
		return nil
	})

	return payload, live, err
}
