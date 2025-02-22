package raftstore

import (
	"bytes"

	"github.com/dgraph-io/badger/v4"

	"github.com/hashicorp/go-msgpack/v2/codec"
	"github.com/hashicorp/raft"
)

var PREFIX_LOG = []byte{0x00}

func (s *Store) FirstIndex() (uint64, error) {
	var value uint64

	err := s.conn.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Reverse:        false,
		})
		defer it.Close()

		it.Seek(PREFIX_LOG)
		if it.ValidForPrefix(PREFIX_LOG) {
			value = bytesToUint64(it.Item().Key()[1:])
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s *Store) LastIndex() (uint64, error) {
	var value uint64

	err := s.conn.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Reverse:        true,
		})
		defer it.Close()

		it.Seek(append(PREFIX_LOG, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF))
		if it.ValidForPrefix(PREFIX_LOG) {
			value = bytesToUint64(it.Item().Key()[1:])
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (s *Store) GetLog(index uint64, log *raft.Log) error {
	return s.conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(PREFIX_LOG, uint64ToBytes(index)...))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				return raft.ErrLogNotFound

			default:
				return err
			}
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return decodeMsgPack(val, log)
	})
}

func (s *Store) StoreLog(log *raft.Log) error {
	val, err := encodeMsgPack(log)
	if err != nil {
		return err
	}

	return s.conn.Update(func(txn *badger.Txn) error {
		return txn.Set(append(PREFIX_LOG, uint64ToBytes(log.Index)...), val)
	})
}

func (s *Store) StoreLogs(logs []*raft.Log) error {
	txn := s.conn.NewTransaction(true)

	for i, log := range logs {
		key := append(PREFIX_LOG, uint64ToBytes(log.Index)...)
		val, err := encodeMsgPack(log)
		if err != nil {
			return err
		}

		if err := txn.Set(key, val); err != nil {
			if err == badger.ErrTxnTooBig {
				err = txn.Commit()
				if err != nil {
					return err
				}

				return s.StoreLogs(logs[i:])
			}

			return err
		}
	}

	return txn.Commit()
}

func (s *Store) DeleteRange(min uint64, max uint64) error {
	txn := s.conn.NewTransaction(true)
	it := txn.NewIterator(badger.IteratorOptions{
		PrefetchValues: false,
		Reverse:        false,
	})

	start := append(PREFIX_LOG, uint64ToBytes(min)...)
	for it.Seek(start); it.Valid(); it.Next() {
		key := it.Item().KeyCopy(nil)

		if bytesToUint64(key[1:]) > max {
			break
		}

		if err := txn.Delete(key); err != nil {
			if err == badger.ErrTxnTooBig {
				it.Close()
				err = txn.Commit()
				if err != nil {
					return err
				}

				return s.DeleteRange(bytesToUint64(key[1:]), max)
			}

			it.Close()
			return err
		}
	}

	it.Close()
	return txn.Commit()
}

func encodeMsgPack(log *raft.Log) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	if err := codec.NewEncoder(buf, &codec.MsgpackHandle{}).Encode(log); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeMsgPack(buf []byte, log *raft.Log) error {
	return codec.
		NewDecoder(bytes.NewBuffer(buf), &codec.MsgpackHandle{}).
		Decode(log)
}
