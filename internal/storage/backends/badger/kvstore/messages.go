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
	w     io.Writer
	since uint64
}

type restoreOperation struct {
	r io.Reader
}

type viewOperation struct {
	txnFn func(txn *badger.Txn) error
}

type updateOperation struct {
	txnFn func(txn *badger.Txn) error
}

var _ operation = (*backupOperation)(nil)
var _ operation = (*restoreOperation)(nil)
var _ operation = (*viewOperation)(nil)
var _ operation = (*updateOperation)(nil)

func (m *backupOperation) Handle(db *badger.DB) error {
	var err error
	stream := db.NewStream()
	stream.NumGo = 1
	stream.LogPrefix = "DB.Backup"
	stream.SinceTs = m.since
	m.since, err = stream.Backup(m.w, m.since)
	return err
}

func (m *restoreOperation) Handle(db *badger.DB) error {
	return db.Load(m.r, 1)
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
