package auth

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"link-society.com/flowg/internal/hash"
)

type UserSystem struct {
	backend *Database
}

func NewUserSystem(backend *Database) *UserSystem {
	return &UserSystem{backend: backend}
}

func (sys *UserSystem) ListUsers() ([]User, error) {
	var users []User = nil

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		usernames, err := sys.fetchUsernames(txn)
		if err != nil {
			return err
		}

		users = []User{}

		for _, username := range usernames {
			user, err := sys.fetchUser(txn, username)
			if err != nil {
				return err
			}

			if user != nil {
				users = append(users, *user)
			}
		}

		return nil
	})

	return users, err
}

func (sys *UserSystem) GetUser(name string) (*User, error) {
	var user *User = nil

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		var err error
		user, err = sys.fetchUser(txn, name)
		return err
	})

	return user, err
}

func (sys *UserSystem) SaveUser(user User, password string) error {
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return sys.backend.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("index:user:%s", user.Name))
		err := txn.Set(key, []byte{})
		if err != nil {
			return fmt.Errorf("failed to save index of user '%s': %w", user.Name, err)
		}

		key = []byte(fmt.Sprintf("user:%s:password", user.Name))
		err = txn.Set(key, []byte(passwordHash))
		if err != nil {
			return fmt.Errorf("failed to save password of user '%s': %w", user.Name, err)
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
					return fmt.Errorf(
						"failed to add role '%s' to user '%s': %w",
						role, user.Name, err,
					)
				}
			} else {
				delete(obsoleteRoles, role)
			}
		}

		for role, key := range obsoleteRoles {
			err := txn.Delete(key)
			if err != nil && err != badger.ErrKeyNotFound {
				return fmt.Errorf(
					"failed to delete role '%s' from user '%s': %w",
					role, user.Name, err,
				)
			}
		}

		return nil
	})
}

func (sys *UserSystem) DeleteUser(name string) error {
	return sys.backend.db.Update(func(txn *badger.Txn) error {
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
				return fmt.Errorf(
					"failed to delete key '%s' of user '%s': %w",
					name, key, err,
				)
			}
		}

		err := txn.Delete([]byte(fmt.Sprintf("index:user:%s", name)))
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf(
				"failed to delete index of user '%s': %w",
				name, err,
			)
		}

		return nil
	})
}

func (sys *UserSystem) VerifyUserPassword(name, password string) (bool, error) {
	isValid := false

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("user:%s:password", name))
		item, err := txn.Get(key)
		if err != nil {
			return fmt.Errorf("failed to get password of user '%s': %w", name, err)
		}

		err = item.Value(func(val []byte) error {
			passwordHash := string(val)
			isValid, err = hash.VerifyPassword(password, passwordHash)
			if err != nil {
				return fmt.Errorf("failed to verify password of user '%s': %w", name, err)
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

func (sys *UserSystem) VerifyUserPermission(name string, scope Scope) (bool, error) {
	hasPermission := false

	err := sys.backend.db.View(func(txn *badger.Txn) error {
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

				case scope == SCOPE_READ_ACLS && roleScope == SCOPE_WRITE_ACLS:
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

func (sys *UserSystem) ListUserScopes(username string) ([]Scope, error) {
	scopeMap := map[Scope]struct{}{}

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", username))
		it := txn.NewIterator(opts)
		defer it.Close()

		roles := make([]string, 0)

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			roleName := string(key[len(username)+11:])
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

				switch roleScope {
				case SCOPE_WRITE_PIPELINES:
					scopeMap[SCOPE_READ_PIPELINES] = struct{}{}
					scopeMap[SCOPE_WRITE_PIPELINES] = struct{}{}

				case SCOPE_WRITE_TRANSFORMERS:
					scopeMap[SCOPE_READ_TRANSFORMERS] = struct{}{}
					scopeMap[SCOPE_WRITE_TRANSFORMERS] = struct{}{}

				case SCOPE_WRITE_STREAMS:
					scopeMap[SCOPE_READ_STREAMS] = struct{}{}
					scopeMap[SCOPE_WRITE_STREAMS] = struct{}{}

				case SCOPE_WRITE_ACLS:
					scopeMap[SCOPE_READ_ACLS] = struct{}{}
					scopeMap[SCOPE_WRITE_ACLS] = struct{}{}

				default:
					scopeMap[roleScope] = struct{}{}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	scopes := make([]Scope, 0, len(scopeMap))
	for scope := range scopeMap {
		scopes = append(scopes, scope)
	}

	return scopes, nil
}

func (sys *UserSystem) fetchUsernames(txn *badger.Txn) ([]string, error) {
	usernames := []string{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte("index:user:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		usernames = append(usernames, string(key[11:]))
	}

	return usernames, nil
}

func (sys *UserSystem) fetchUser(txn *badger.Txn, name string) (*User, error) {
	_, err := txn.Get([]byte(fmt.Sprintf("index:user:%s", name)))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get index of user '%s': %w", name, err)
	}

	user := &User{Name: name}

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

	return user, nil
}
