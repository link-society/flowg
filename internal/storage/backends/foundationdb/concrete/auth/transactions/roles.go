package transactions

import (
	"fmt"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/models"
)

// Subspace layout for roles:
//
//	<root>/index/role/<name>  → marker (empty)  — lists all roles
//	<root>/role/<name>/<scope> → marker (empty)  — scopes per role
var (
	roleSub      subspace.Subspace // <root>/role
	indexRoleSub subspace.Subspace // <root>/index/role
)

func initRoleSubspaces(root subspace.Subspace) {
	roleSub = root.Sub("role")
	indexRoleSub = root.Sub("index").Sub("role")
}

// Init initializes all package-level subspaces from the given root.
func Init(root subspace.Subspace) {
	initRoleSubspaces(root)
	initUserSubspaces(root)
	initTokenSubspaces(root)
}

// ListRoles loads every role together with the scopes it grants.
func ListRoles(tr fdb.ReadTransaction) ([]models.Role, error) {
	roles := []models.Role{}
	roleNames := fetchRoleNames(tr)

	for _, roleName := range roleNames {
		role, err := FetchRole(tr, roleName)
		if err != nil {
			return nil, err
		}

		if role != nil {
			roles = append(roles, *role)
		}
	}

	return roles, nil
}

// FetchRole reconstructs a single role by collecting its
// <root>/role/<name>/<scope> keys; the key itself carries the scope, so values
// are never read. Returns nil when the role does not exist.
func FetchRole(tr fdb.ReadTransaction, name string) (*models.Role, error) {
	val, err := tr.Get(indexRoleSub.Pack(tuple.Tuple{name})).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get index of role '%s': %w", name, err)
	}
	if val == nil {
		return nil, nil
	}

	role := &models.Role{Name: name}
	sub := roleSub.Sub(name)

	iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := sub.Unpack(kv.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack scope key for role '%s': %w", name, err)
		}

		scopeStr := t[0].(string)
		scope, err := models.ParseScope(scopeStr)
		if err != nil {
			return nil, err
		}

		role.Scopes = append(role.Scopes, scope)
	}

	return role, nil
}

// SaveRole persists a role and reconciles its scope keys against what is already
// stored: it writes the index marker, then converges the existing scope keys
// towards the desired scope set by adding the missing ones and deleting the
// obsolete ones.
func SaveRole(tr fdb.Transaction, role models.Role) error {
	key := indexRoleSub.Pack(tuple.Tuple{role.Name})
	tr.Set(key, []byte{})

	sub := roleSub.Sub(role.Name)

	obsoleteScopes := map[models.Scope]struct{}{}

	// Collect the scopes currently stored for the role, tentatively flagging any
	// that are no longer part of the desired set as obsolete.
	iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := sub.Unpack(kv.Key)
		if err != nil {
			return fmt.Errorf("failed to unpack scope key for role '%s': %w", role.Name, err)
		}

		scopeStr := t[0].(string)
		scope, err := models.ParseScope(scopeStr)
		if err != nil {
			return err
		}

		if !role.HasScope(scope) {
			obsoleteScopes[scope] = struct{}{}
		}
	}

	// Persist the desired scopes: create the ones that are missing and clear the
	// obsolete flag on the ones that already exist.
	for _, scope := range role.Scopes {
		if _, exists := obsoleteScopes[scope]; !exists {
			tr.Set(sub.Pack(tuple.Tuple{string(scope)}), []byte{})
		} else {
			delete(obsoleteScopes, scope)
		}
	}

	// Whatever scope keys remain flagged are no longer granted and get removed.
	for scope := range obsoleteScopes {
		tr.Clear(sub.Pack(tuple.Tuple{string(scope)}))
	}

	return nil
}

// DeleteRole removes a role entirely: every <root>/role/<name>/<scope> key plus
// its index marker.
func DeleteRole(tr fdb.Transaction, name string) error {
	sub := roleSub.Sub(name)
	tr.ClearRange(sub)
	tr.Clear(indexRoleSub.Pack(tuple.Tuple{name}))

	return nil
}

// fetchRoleNames lists existing role names from their index markers.
func fetchRoleNames(tr fdb.ReadTransaction) []string {
	roleNames := []string{}

	iter := tr.GetRange(indexRoleSub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := indexRoleSub.Unpack(kv.Key)
		if err != nil {
			continue
		}

		roleNames = append(roleNames, t[0].(string))
	}

	return roleNames
}
