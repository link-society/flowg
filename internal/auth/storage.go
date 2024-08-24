package auth

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/logging"
)

type Database struct {
	db *badger.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	opts := badger.
		DefaultOptions(dbPath).
		WithLogger(&logging.BadgerLogger{Channel: "authdb"}).
		WithCompression(options.ZSTD)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) SaveRole(role Role) error {
	return d.db.Update(func(txn *badger.Txn) error {
		for _, scope := range role.Scopes {
			key := []byte(fmt.Sprintf("role:%s:%s", role.Name, scope))
			err := txn.Set(key, []byte{})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *Database) DeleteRole(name string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("role:%s:", name))
		it := txn.NewIterator(opts)
		defer it.Close()

		keys := make([][]byte, 0)

		for it.Rewind(); it.Valid(); it.Next() {
			keys = append(keys, it.Item().KeyCopy(nil))
		}

		for _, key := range keys {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
