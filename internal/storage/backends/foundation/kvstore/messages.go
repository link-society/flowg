package kvstore

import "github.com/apple/foundationdb/bindings/go/src/fdb"

type message struct {
	replyTo chan<- error
	operation
}

type operation interface {
	Handle(db fdb.Database) error
}

type viewOperation struct {
	txnFn func(txn fdb.ReadTransaction) error
}

type updateOperation struct {
	txnFn func(txn fdb.Transaction) error
}

var _ operation = (*viewOperation)(nil)
var _ operation = (*updateOperation)(nil)

func (m *viewOperation) Handle(db fdb.Database) error {
	_, err := db.ReadTransact(func(txn fdb.ReadTransaction) (any, error) {
		return nil, m.txnFn(txn)
	})

	return err
}

func (m *updateOperation) Handle(db fdb.Database) error {
	_, err := db.Transact(func(txn fdb.Transaction) (any, error) {
		return nil, m.txnFn(txn)
	})
	return err
}
