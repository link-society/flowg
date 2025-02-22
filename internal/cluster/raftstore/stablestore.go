package raftstore

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

var PREFIX_STABLE = []byte{0x01}

func (s *Store) Get(key []byte) ([]byte, error) {
	var value []byte

	err := s.conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get(append(PREFIX_STABLE, key...))
		if err != nil {
			switch {
			case err == badger.ErrKeyNotFound:
				return fmt.Errorf("not found")

			default:
				return err
			}
		}

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *Store) Set(key []byte, val []byte) error {
	return s.conn.Update(func(txn *badger.Txn) error {
		return txn.Set(append(PREFIX_STABLE, key...), val)
	})
}

func (s *Store) GetUint64(key []byte) (uint64, error) {
	value, err := s.Get(key)
	if err != nil {
		return 0, err
	}

	return bytesToUint64(value), nil
}

func (s *Store) SetUint64(key []byte, val uint64) error {
	return s.Set(key, uint64ToBytes(val))
}
