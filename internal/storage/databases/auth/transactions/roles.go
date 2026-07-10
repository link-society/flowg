package transactions

import (
	"fmt"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

// ListRoles loads every role together with the scopes it grants.
func ListRoles(txn kv.QueryTx) ([]models.Role, error) {
	roles := []models.Role{}
	roleNames := fetchRoleNames(txn)

	for _, roleName := range roleNames {
		role, err := FetchRole(txn, roleName)
		if err != nil {
			return nil, fmt.Errorf("failed to list roles: %w", err)
		}

		if role != nil {
			roles = append(roles, *role)
		}
	}

	return roles, nil
}

// FetchRole reconstructs a single role by collecting its "role:<name>:*" scope
// keys; the key itself carries the scope, so values are never read.
func FetchRole(txn kv.QueryTx, name string) (*models.Role, error) {
	role := &models.Role{Name: name}

	prefix := kv.Key{"role", name}

	for key := range txn.IterKeys(prefix, kv.KeyRange{}) {
		scopeName := key[len(key)-1]
		scope, err := models.ParseScope(scopeName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse scope %q while fetching role %q: %w", scopeName, name, err)
		}

		role.Scopes = append(role.Scopes, scope)
	}

	return role, nil
}

// SaveRole persists a role and reconciles its scope keys against what is already
// stored: it writes the "index:role:<name>" marker, then converges the existing
// "role:<name>:*" keys towards the desired scope set by adding the missing ones
// and deleting the obsolete ones.
func SaveRole(txn kv.MutationTx, role models.Role) error {
	key := kv.Key{"index", "role", role.Name}
	err := txn.Set(key, []byte{})
	if err != nil {
		return fmt.Errorf("failed to write index for role %q: %w", role.Name, err)
	}

	obsoleteScopes := map[models.Scope]kv.Key{}

	// Collect the scopes currently stored for the role, tentatively flagging any
	// that are no longer part of the desired set as obsolete.
	for key := range txn.IterKeys(kv.Key{"role", role.Name}, kv.KeyRange{}) {
		scopeName := key[len(key)-1]
		scope, err := models.ParseScope(scopeName)
		if err != nil {
			return fmt.Errorf("failed to parse scope %q while saving role %q: %w", scopeName, role.Name, err)
		}

		if !role.HasScope(scope) {
			obsoleteScopes[scope] = key
		}
	}

	// Persist the desired scopes: create the ones that are missing and clear the
	// obsolete flag on the ones that already exist.
	for _, scope := range role.Scopes {
		if _, exists := obsoleteScopes[scope]; !exists {
			key := kv.Key{"role", role.Name, string(scope)}
			err := txn.Set(key, []byte{})
			if err != nil {
				return fmt.Errorf("failed to write scope %q for role %q: %w", scope, role.Name, err)
			}
		} else {
			delete(obsoleteScopes, scope)
		}
	}

	// Whatever scope keys remain flagged are no longer granted and get removed.
	for _, key := range obsoleteScopes {
		err := txn.Clear(key)
		if err != nil {
			return fmt.Errorf("failed to clear scope %q for role %q: %w", key[len(key)-1], role.Name, err)
		}
	}

	return nil
}

// DeleteRole removes a role entirely: every "role:<name>:*" scope key plus its
// "index:role:<name>" marker.
func DeleteRole(txn kv.MutationTx, name string) error {
	keys := make([]kv.Key, 0)

	for key := range txn.IterKeys(kv.Key{"role", name}, kv.KeyRange{}) {
		keys = append(keys, key)
	}

	for _, key := range keys {
		if err := txn.Clear(key); err != nil {
			return fmt.Errorf("failed to clear scope %q for role %q: %w", key[len(key)-1], name, err)
		}
	}

	if err := txn.Clear(kv.Key{"index", "role", name}); err != nil {
		return fmt.Errorf("failed to clear index for role %q: %w", name, err)
	}

	return nil
}

// fetchRoleNames lists existing role names from their "index:role:" markers.
func fetchRoleNames(txn kv.QueryTx) []string {
	roleNames := []string{}

	for key := range txn.IterKeys(kv.Key{"index", "role"}, kv.KeyRange{}) {
		roleName := key[len(key)-1]
		roleNames = append(roleNames, roleName)
	}

	return roleNames
}
