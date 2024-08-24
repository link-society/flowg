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

func (d *Database) ListRoles() ([]string, error) {
	roles := make([]string, 0)

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("index:role:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			roles = append(roles, string(key[11:]))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (d *Database) GetRole(name string) (Role, error) {
	role := Role{Name: name}

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("role:%s:", name))
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			scopeName := string(key[len(name)+6:])
			scope, err := ParseScope(scopeName)
			if err != nil {
				return err
			}

			role.Scopes = append(role.Scopes, scope)
		}

		return nil
	})

	return role, err
}

func (d *Database) SaveRole(role Role) error {
	return d.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("index:role:%s", role.Name))
		err := txn.Set(key, []byte{})
		if err != nil {
			return err
		}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("role:%s:", role.Name))
		it := txn.NewIterator(opts)
		defer it.Close()

		obsoleteScopes := map[Scope][]byte{}

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)
			scopeName := string(key[len(role.Name)+6:])
			scope, err := ParseScope(scopeName)
			if err != nil {
				return err
			}

			if !role.HasScope(scope) {
				obsoleteScopes[scope] = key
			}
		}

		for _, key := range obsoleteScopes {
			err := txn.Delete(key)
			if err != nil {
				return err
			}
		}

		for _, scope := range role.Scopes {
			if _, exists := obsoleteScopes[scope]; !exists {
				key := []byte(fmt.Sprintf("role:%s:%s", role.Name, scope))
				err := txn.Set(key, []byte{})
				if err != nil {
					return err
				}
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
			if err != nil && err != badger.ErrKeyNotFound {
				return err
			}
		}

		err := txn.Delete([]byte(fmt.Sprintf("index:role:%s", name)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		return nil
	})
}
