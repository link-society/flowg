package transactions

import (
	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

func setItem(txn *badger.Txn, key []byte, payload []byte, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	env := lww.Envelope{Timestamp: ts, Payload: payload}
	applied, err := lww.Apply(txn, key, env)
	if err != nil {
		return err
	}
	if applied && sink != nil {
		*sink = append(*sink, changefeed.Record{
			Key:   append([]byte(nil), key...),
			Value: env.Marshal(),
		})
	}
	return nil
}

func deleteItem(txn *badger.Txn, key []byte, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	env := lww.Envelope{Timestamp: ts, Deleted: true}
	applied, err := lww.Apply(txn, key, env)
	if err != nil {
		return err
	}
	if applied && sink != nil {
		*sink = append(*sink, changefeed.Record{
			Key:   append([]byte(nil), key...),
			Value: env.Marshal(),
		})
	}
	return nil
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
