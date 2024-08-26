package auth

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"github.com/google/uuid"

	"link-society.com/flowg/internal/hash"
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
			return &PersistError{Operation: "add-role-index", Reason: err}
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

		for _, scope := range role.Scopes {
			if _, exists := obsoleteScopes[scope]; !exists {
				key := []byte(fmt.Sprintf("role:%s:%s", role.Name, scope))
				err := txn.Set(key, []byte{})
				if err != nil {
					return &PersistError{Operation: "add-role-scope", Reason: err}
				}
			} else {
				delete(obsoleteScopes, scope)
			}
		}

		for _, key := range obsoleteScopes {
			err := txn.Delete(key)
			if err != nil {
				return &PersistError{Operation: "delete-role-obsolete-scope", Reason: err}
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
				return &PersistError{Operation: "delete-role-scope", Reason: err}
			}
		}

		err := txn.Delete([]byte(fmt.Sprintf("index:role:%s", name)))
		if err != nil && err != badger.ErrKeyNotFound {
			return &PersistError{Operation: "delete-role-index", Reason: err}
		}

		return nil
	})
}

func (d *Database) ListUsers() ([]string, error) {
	users := make([]string, 0)

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte("index:user:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			users = append(users, string(key[11:]))
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (d *Database) GetUser(name string) (*User, error) {
	var user *User = nil

	err := d.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(fmt.Sprintf("index:user:%s", name)))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}

			return &QueryError{Operation: "get-user-index", Reason: err}
		}

		user = &User{Name: name}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", name))
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			roleName := string(key[len(name)+11:])
			user.Roles = append(user.Roles, roleName)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Database) SaveUser(user User, password string) error {
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return &PersistError{Operation: "hash-password", Reason: err}
	}

	return d.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("index:user:%s", user.Name))
		err := txn.Set(key, []byte{})
		if err != nil {
			return &PersistError{Operation: "add-user-index", Reason: err}
		}

		key = []byte(fmt.Sprintf("user:%s:password", user.Name))
		err = txn.Set(key, []byte(passwordHash))
		if err != nil {
			return &PersistError{Operation: "add-user-password", Reason: err}
		}

		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", user.Name))
		it := txn.NewIterator(opts)
		defer it.Close()

		obsoleteRoles := map[string][]byte{}

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)
			roleName := string(key[len(user.Name)+11:])
			if !user.HasRole(roleName) {
				obsoleteRoles[roleName] = key
			}
		}

		for _, role := range user.Roles {
			if _, exists := obsoleteRoles[role]; !exists {
				key := []byte(fmt.Sprintf("user:%s:role:%s", user.Name, role))
				err := txn.Set(key, []byte{})
				if err != nil {
					return &PersistError{Operation: "add-user-role", Reason: err}
				}
			} else {
				delete(obsoleteRoles, role)
			}
		}

		for _, key := range obsoleteRoles {
			err := txn.Delete(key)
			if err != nil && err != badger.ErrKeyNotFound {
				return &PersistError{Operation: "delete-user-obsolete-role", Reason: err}
			}
		}

		return nil
	})
}

func (d *Database) DeleteUser(name string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		keys := make([][]byte, 0)

		(func() {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(fmt.Sprintf("user:%s:", name))
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				keys = append(keys, it.Item().KeyCopy(nil))
			}
		})()

		(func() {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(fmt.Sprintf("pat:%s:", name))
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				keys = append(keys, it.Item().KeyCopy(nil))
			}
		})()

		for _, key := range keys {
			err := txn.Delete(key)
			if err != nil && err != badger.ErrKeyNotFound {
				return &PersistError{Operation: "delete-user-property", Reason: err}
			}
		}

		err := txn.Delete([]byte(fmt.Sprintf("index:user:%s", name)))
		if err != nil && err != badger.ErrKeyNotFound {
			return &PersistError{Operation: "delete-user-index", Reason: err}
		}

		return nil
	})
}

func (d *Database) VerifyUserPassword(name, password string) (bool, error) {
	isValid := false

	err := d.db.View(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("user:%s:password", name))
		item, err := txn.Get(key)
		if err != nil {
			return &QueryError{Operation: "user-password", Reason: err}
		}

		err = item.Value(func(val []byte) error {
			passwordHash := string(val)
			isValid, err = hash.VerifyPassword(password, passwordHash)
			if err != nil {
				return &QueryError{Operation: "verify-user-password", Reason: err}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return isValid, nil
}

func (d *Database) VerifyUserPermission(name string, scope Scope) (bool, error) {
	hasPermission := false

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", name))
		it := txn.NewIterator(opts)
		defer it.Close()

		roles := make([]string, 0)

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			roleName := string(key[len(name)+11:])
			roles = append(roles, roleName)
		}

		for _, roleName := range roles {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			opts.Prefix = []byte(fmt.Sprintf("role:%s:", roleName))
			it := txn.NewIterator(opts)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				key := it.Item().Key()
				scopeName := string(key[len(roleName)+6:])
				roleScope, err := ParseScope(scopeName)
				if err != nil {
					return err
				}

				switch {
				case scope == roleScope:
					hasPermission = true
					return nil

				case scope == SCOPE_READ_PIPELINES && roleScope == SCOPE_WRITE_PIPELINES:
					hasPermission = true
					return nil

				case scope == SCOPE_READ_TRANSFORMERS && roleScope == SCOPE_WRITE_TRANSFORMERS:
					hasPermission = true
					return nil

				case scope == SCOPE_READ_STREAMS && roleScope == SCOPE_WRITE_STREAMS:
					hasPermission = true
					return nil
				}
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

func (d *Database) AddPersonalAccessToken(username string, token string) error {
	tokenHash, err := hash.HashPassword(token)
	if err != nil {
		return &PersistError{Operation: "hash-token", Reason: err}
	}

	return d.db.Update(func(txn *badger.Txn) error {
		userKey := []byte(fmt.Sprintf("index:user:%s", username))
		_, err := txn.Get(userKey)
		if err != nil {
			return &QueryError{Operation: "get-user", Reason: err}
		}

		tokenKey := []byte(fmt.Sprintf("pat:%s:%s", username, uuid.New().String()))
		err = txn.Set(tokenKey, []byte(tokenHash))
		if err != nil {
			return &PersistError{Operation: "add-token", Reason: err}
		}

		return nil
	})
}

func (d *Database) VerifyPersonalAccessToken(token string) (string, bool, error) {
	var username string
	found := false

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("pat:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			keySuffix := string(key[4:])
			associatedUser := keySuffix[:len(keySuffix)-37] // remove UUID

			err := it.Item().Value(func(val []byte) error {
				tokenHash := string(val)
				isValid, err := hash.VerifyPassword(token, tokenHash)
				if err != nil {
					return &QueryError{Operation: "verify-token", Reason: err}
				}

				if isValid {
					username = associatedUser
					found = true
				}

				return nil
			})

			if err != nil {
				return err
			}

			if found {
				break
			}
		}

		return nil
	})

	if err != nil {
		return "", false, err
	}

	return username, found, nil
}
