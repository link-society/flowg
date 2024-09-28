package kvstore

import (
	"github.com/dgraph-io/badger/v3"
)

type message struct {
	replyTo chan<- error
	operation
}

type operation interface {
	Handle(db *badger.DB) error
}

type viewOperation struct {
	txnFn func(txn *badger.Txn) error
}

type updateOperation struct {
	txnFn func(txn *badger.Txn) error
}

func (m *viewOperation) Handle(db *badger.DB) error {
	return db.View(m.txnFn)
}

func (m *updateOperation) Handle(db *badger.DB) error {
	for {
		err := db.Update(m.txnFn)

		switch err {
		case nil:
			return nil

		case badger.ErrConflict:
			db.Opts().Logger.Debugf("Conflict detected, retrying transaction")
			continue

		default:
			return err
		}
	}
}
