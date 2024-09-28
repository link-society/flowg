package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"

	"link-society.com/flowg/internal/models"
)

func ListRoles(txn *badger.Txn) ([]models.Role, error) {
	roles := []models.Role{}
	roleNames := fetchRoleNames(txn)

	for _, roleName := range roleNames {
		role, err := FetchRole(txn, roleName)
		if err != nil {
			return nil, err
		}

		if role != nil {
			roles = append(roles, *role)
		}
	}

	return roles, nil
}

func FetchRole(txn *badger.Txn, name string) (*models.Role, error) {
	role := &models.Role{Name: name}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("role:%s:", name))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().Key()
		scopeName := string(key[len(name)+6:])
		scope, err := models.ParseScope(scopeName)
		if err != nil {
			return nil, err
		}

		role.Scopes = append(role.Scopes, scope)
	}

	return role, nil
}

func SaveRole(txn *badger.Txn, role models.Role) error {
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

	obsoleteScopes := map[models.Scope][]byte{}

	for it.Rewind(); it.Valid(); it.Next() {
		key := it.Item().KeyCopy(nil)
		scopeName := string(key[len(role.Name)+6:])
		scope, err := models.ParseScope(scopeName)
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
}

func DeleteRole(txn *badger.Txn, name string) error {
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
}

func fetchRoleNames(txn *badger.Txn) []string {
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
