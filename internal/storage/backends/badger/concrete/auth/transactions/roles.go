package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
)

// ListRoles loads every role together with the scopes it grants.
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

// FetchRole reconstructs a single role by collecting its "role:<name>:*" scope
// keys; the key itself carries the scope, so values are never read.
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

// SaveRole persists a role and reconciles its scope keys against what is already
// stored: it writes the "index:role:<name>" marker, then converges the existing
// "role:<name>:*" keys towards the desired scope set by adding the missing ones
// and deleting the obsolete ones.
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

	// Collect the scopes currently stored for the role, tentatively flagging any
	// that are no longer part of the desired set as obsolete.
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

	// Persist the desired scopes: create the ones that are missing and clear the
	// obsolete flag on the ones that already exist.
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

	// Whatever scope keys remain flagged are no longer granted and get removed.
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

// DeleteRole removes a role entirely: every "role:<name>:*" scope key plus its
// "index:role:<name>" marker.
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

// fetchRoleNames lists existing role names from their "index:role:" markers.
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
