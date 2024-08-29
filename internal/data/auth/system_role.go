package auth

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type RoleSystem struct {
	backend *Database
}

func NewRoleSystem(backend *Database) *RoleSystem {
	return &RoleSystem{backend: backend}
}

func (sys *RoleSystem) ListRoles() ([]Role, error) {
	roles := []Role{}

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		roleNames := sys.fetchRoleNames(txn)

		for _, roleName := range roleNames {
			role, err := sys.fetchRole(txn, roleName)
			if err != nil {
				return err
			}

			if role != nil {
				roles = append(roles, *role)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (sys *RoleSystem) GetRole(name string) (*Role, error) {
	var role *Role

	err := sys.backend.db.View(func(txn *badger.Txn) error {
		var err error
		role, err = sys.fetchRole(txn, name)
		return err
	})

	return role, err
}

func (sys *RoleSystem) SaveRole(role Role) error {
	return sys.backend.db.Update(func(txn *badger.Txn) error {
		key := []byte(fmt.Sprintf("index:role:%s", role.Name))
		err := txn.Set(key, []byte{})
		if err != nil {
			return fmt.Errorf(
				"failed to save index of role '%s': %w",
				role.Name, err,
			)
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
					return fmt.Errorf(
						"failed to add scope '%s' to role '%s': %w",
						scope, role.Name, err,
					)
				}
			} else {
				delete(obsoleteScopes, scope)
			}
		}

		for scope, key := range obsoleteScopes {
			err := txn.Delete(key)
			if err != nil {
				return fmt.Errorf(
					"failed to remove obsolete scope '%s' from role '%s': %w",
					scope, role.Name, err,
				)
			}
		}

		return nil
	})
}

func (sys *RoleSystem) DeleteRole(name string) error {
	return sys.backend.db.Update(func(txn *badger.Txn) error {
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
				return fmt.Errorf(
					"failed to delete key '%s' from role '%s': %w",
					string(key), name, err,
				)
			}
		}

		err := txn.Delete([]byte(fmt.Sprintf("index:role:%s", name)))
		if err != nil && err != badger.ErrKeyNotFound {
			return fmt.Errorf(
				"failed to delete index of role '%s': %w",
				name, err,
			)
		}

		return nil
	})
}

func (sys *RoleSystem) fetchRoleNames(txn *badger.Txn) []string {
	roleNames := []string{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte("index:role:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		roleNames = append(roleNames, string(key[11:]))
	}

	return roleNames
}

func (sys *RoleSystem) fetchRole(txn *badger.Txn, name string) (*Role, error) {
	role := &Role{Name: name}

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
			return nil, err
		}

		role.Scopes = append(role.Scopes, scope)
	}

	return role, nil
}
