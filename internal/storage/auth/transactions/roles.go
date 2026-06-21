package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/utils/hlc"
)

func ListRoles(txn *badger.Txn) ([]models.Role, error) {
	roles := []models.Role{}
	roleNames, err := fetchRoleNames(txn)
	if err != nil {
		return nil, err
	}

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
	opts.Prefix = []byte(fmt.Sprintf("role:%s:", name))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			return nil, err
		}
		if !live {
			continue
		}

		scopeName := string(item.Key()[len(name)+6:])
		scope, err := models.ParseScope(scopeName)
		if err != nil {
			return nil, err
		}

		role.Scopes = append(role.Scopes, scope)
	}

	return role, nil
}

func SaveRole(txn *badger.Txn, role models.Role, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	indexKey := []byte(fmt.Sprintf("index:role:%s", role.Name))
	if err := setItem(txn, indexKey, []byte{}, ts, sink); err != nil {
		return fmt.Errorf(
			"failed to save index of role '%s': %w",
			role.Name, err,
		)
	}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("role:%s:", role.Name))
	it := txn.NewIterator(opts)

	obsoleteScopes := map[models.Scope][]byte{}

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			it.Close()
			return err
		}
		if !live {
			continue
		}

		key := item.KeyCopy(nil)
		scopeName := string(key[len(role.Name)+6:])
		scope, err := models.ParseScope(scopeName)
		if err != nil {
			it.Close()
			return err
		}

		if !role.HasScope(scope) {
			obsoleteScopes[scope] = key
		}
	}
	it.Close()

	for _, scope := range role.Scopes {
		if _, exists := obsoleteScopes[scope]; !exists {
			key := []byte(fmt.Sprintf("role:%s:%s", role.Name, scope))
			if err := setItem(txn, key, []byte{}, ts, sink); err != nil {
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
		if err := deleteItem(txn, key, ts, sink); err != nil {
			return fmt.Errorf(
				"failed to remove obsolete scope '%s' from role '%s': %w",
				scope, role.Name, err,
			)
		}
	}

	return nil
}

func DeleteRole(txn *badger.Txn, name string, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("role:%s:", name))
	it := txn.NewIterator(opts)

	keys := make([][]byte, 0)

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			it.Close()
			return err
		}
		if live {
			keys = append(keys, item.KeyCopy(nil))
		}
	}
	it.Close()

	for _, key := range keys {
		if err := deleteItem(txn, key, ts, sink); err != nil {
			return fmt.Errorf(
				"failed to delete key '%s' from role '%s': %w",
				string(key), name, err,
			)
		}
	}

	indexKey := []byte(fmt.Sprintf("index:role:%s", name))
	if err := deleteItem(txn, indexKey, ts, sink); err != nil {
		return fmt.Errorf(
			"failed to delete index of role '%s': %w",
			name, err,
		)
	}

	return nil
}

func fetchRoleNames(txn *badger.Txn) ([]string, error) {
	roleNames := []string{}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte("index:role:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			return nil, err
		}
		if !live {
			continue
		}
		roleNames = append(roleNames, string(item.Key()[11:]))
	}

	return roleNames, nil
}
