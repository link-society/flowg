package kvstore

import (
	"io"

	"github.com/dgraph-io/badger/v4"
)

type message struct {
	replyTo chan<- error
	operation
}

type operation interface {
	Handle(db *badger.DB) error
}

type backupOperation struct {
	w io.Writer
}

type viewOperation struct {
	txnFn func(txn *badger.Txn) error
}

type updateOperation struct {
	txnFn func(txn *badger.Txn) error
}

func (m *backupOperation) Handle(db *badger.DB) error {
	_, err := db.Backup(m.w, 0)
	return err
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
