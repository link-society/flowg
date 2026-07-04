package transactions

import (
	"fmt"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth/hash"
)

// Subspace layout for users:
//
//	<root>/index/user/<name>          → marker (empty)      — lists all users
//	<root>/user/<name>/"password"      → password hash        — login credentials
//	<root>/user/<name>/role/<role>     → marker (empty)      — role membership
var (
	userSub      subspace.Subspace // <root>/user
	indexUserSub subspace.Subspace // <root>/index/user
)

func initUserSubspaces(root subspace.Subspace) {
	userSub = root.Sub("user")
	indexUserSub = root.Sub("index").Sub("user")
}

// ListUsers loads every user together with their assigned roles.
func ListUsers(tr fdb.ReadTransaction) ([]models.User, error) {
	usernames, err := fetchUsernames(tr)
	if err != nil {
		return nil, err
	}

	users := []models.User{}

	for _, username := range usernames {
		user, err := FetchUser(tr, username)
		if err != nil {
			return nil, err
		}

		if user != nil {
			users = append(users, *user)
		}
	}

	return users, nil
}

// FetchUser loads a single user, returning nil when no index marker exists.
// The user's roles are gathered from their <root>/user/<name>/role/* keys.
func FetchUser(tr fdb.ReadTransaction, name string) (*models.User, error) {
	val, err := tr.Get(indexUserSub.Pack(tuple.Tuple{name})).Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get index of user '%s': %w", name, err)
	}
	if val == nil {
		return nil, nil
	}

	user := &models.User{Name: name}
	roleSub := userSub.Sub(name).Sub("role")

	iter := tr.GetRange(roleSub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := roleSub.Unpack(kv.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack role key for user '%s': %w", name, err)
		}

		user.Roles = append(user.Roles, t[0].(string))
	}

	return user, nil
}

// SaveUser stores the user's password hash under <root>/user/<name>/"password"
// and then reconciles their role assignments via PatchUserRoles.
func SaveUser(tr fdb.Transaction, user models.User, password string) error {
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	key := userSub.Sub(user.Name).Pack(tuple.Tuple{"password"})
	tr.Set(key, []byte(passwordHash))

	return PatchUserRoles(tr, user)
}

// PatchUserRoles writes the index marker and converges the stored role keys
// towards the user's desired role set, adding the missing assignments and
// deleting the obsolete ones (the password is left untouched).
func PatchUserRoles(tr fdb.Transaction, user models.User) error {
	tr.Set(indexUserSub.Pack(tuple.Tuple{user.Name}), []byte{})

	sub := userSub.Sub(user.Name).Sub("role")

	obsoleteRoles := map[string]struct{}{}

	// Collect currently assigned roles, flagging any not in the desired set as
	// obsolete.
	iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := sub.Unpack(kv.Key)
		if err != nil {
			return fmt.Errorf("failed to unpack role key for user '%s': %w", user.Name, err)
		}

		roleName := t[0].(string)
		if !user.HasRole(roleName) {
			obsoleteRoles[roleName] = struct{}{}
		}
	}

	// Create missing roles, clear obsolete flag on existing ones.
	for _, role := range user.Roles {
		if _, exists := obsoleteRoles[role]; !exists {
			tr.Set(sub.Pack(tuple.Tuple{role}), []byte{})
		} else {
			delete(obsoleteRoles, role)
		}
	}

	// Remove whatever roles remain flagged as obsolete.
	for role := range obsoleteRoles {
		tr.Clear(sub.Pack(tuple.Tuple{role}))
	}

	return nil
}

// DeleteUser removes everything tied to a user: all <root>/user/<name>/* keys
// (roles and password), all of their <root>/pat/<name>/* tokens, and the
// index marker.
func DeleteUser(tr fdb.Transaction, name string) error {
	userRange := userSub.Sub(name)
	tr.ClearRange(userRange)

	patRange := patSub.Sub(name)
	tr.ClearRange(patRange)

	tr.Clear(indexUserSub.Pack(tuple.Tuple{name}))

	return nil
}

// VerifyUserPassword argon2id-compares a candidate password against the hash
// stored under <root>/user/<name>/"password". Returns an error when the user
// does not exist or the stored hash is corrupt.
func VerifyUserPassword(tr fdb.ReadTransaction, name, password string) (bool, error) {
	key := userSub.Sub(name).Pack(tuple.Tuple{"password"})
	val, err := tr.Get(key).Get()
	if err != nil {
		return false, fmt.Errorf("failed to get password of user '%s': %w", name, err)
	}
	if val == nil {
		return false, fmt.Errorf("user '%s' not found", name)
	}

	passwordHash := string(val)
	return hash.VerifyPassword(password, passwordHash)
}

// VerifyUserPermission reports whether a user holds a given scope. It walks the
// user's roles and, for each, the scopes it grants, returning true on a direct
// match or when the user holds the matching write scope for a requested read
// scope (write implies read).
func VerifyUserPermission(tr fdb.ReadTransaction, name string, scope models.Scope) (bool, error) {
	user, err := FetchUser(tr, name)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	for _, roleName := range user.Roles {
		sub := roleSub.Sub(roleName)
		iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
		for iter.Advance() {
			kv := iter.MustGet()

			t, err := sub.Unpack(kv.Key)
			if err != nil {
				return false, fmt.Errorf("failed to unpack scope key for role '%s': %w", roleName, err)
			}

			scopeStr := t[0].(string)
			roleScope, err := models.ParseScope(scopeStr)
			if err != nil {
				return false, err
			}

			switch {
			case scope == roleScope:
				return true, nil

			case scope == models.SCOPE_READ_PIPELINES && roleScope == models.SCOPE_WRITE_PIPELINES:
				return true, nil

			case scope == models.SCOPE_READ_TRANSFORMERS && roleScope == models.SCOPE_WRITE_TRANSFORMERS:
				return true, nil

			case scope == models.SCOPE_READ_STREAMS && roleScope == models.SCOPE_WRITE_STREAMS:
				return true, nil

			case scope == models.SCOPE_READ_FORWARDERS && roleScope == models.SCOPE_WRITE_FORWARDERS:
				return true, nil

			case scope == models.SCOPE_READ_ACLS && roleScope == models.SCOPE_WRITE_ACLS:
				return true, nil
			}
		}
	}

	return false, nil
}

// ListUserScopes returns the de-duplicated set of scopes a user effectively
// holds, resolved through their roles. Each write scope additionally pulls in
// its corresponding read scope.
func ListUserScopes(tr fdb.ReadTransaction, username string) ([]models.Scope, error) {
	scopeMap := map[models.Scope]struct{}{}

	user, err := FetchUser(tr, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return []models.Scope{}, nil
	}

	for _, roleName := range user.Roles {
		sub := roleSub.Sub(roleName)
		iter := tr.GetRange(sub, fdb.RangeOptions{}).Iterator()
		for iter.Advance() {
			kv := iter.MustGet()

			t, err := sub.Unpack(kv.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to unpack scope key for role '%s': %w", roleName, err)
			}

			scopeStr := t[0].(string)
			roleScope, err := models.ParseScope(scopeStr)
			if err != nil {
				return nil, err
			}

			switch roleScope {
			case models.SCOPE_WRITE_PIPELINES:
				scopeMap[models.SCOPE_READ_PIPELINES] = struct{}{}
				scopeMap[models.SCOPE_WRITE_PIPELINES] = struct{}{}

			case models.SCOPE_WRITE_TRANSFORMERS:
				scopeMap[models.SCOPE_READ_TRANSFORMERS] = struct{}{}
				scopeMap[models.SCOPE_WRITE_TRANSFORMERS] = struct{}{}

			case models.SCOPE_WRITE_STREAMS:
				scopeMap[models.SCOPE_READ_STREAMS] = struct{}{}
				scopeMap[models.SCOPE_WRITE_STREAMS] = struct{}{}

			case models.SCOPE_WRITE_FORWARDERS:
				scopeMap[models.SCOPE_READ_FORWARDERS] = struct{}{}
				scopeMap[models.SCOPE_WRITE_FORWARDERS] = struct{}{}

			case models.SCOPE_WRITE_ACLS:
				scopeMap[models.SCOPE_READ_ACLS] = struct{}{}
				scopeMap[models.SCOPE_WRITE_ACLS] = struct{}{}

			default:
				scopeMap[roleScope] = struct{}{}
			}
		}
	}

	scopes := make([]models.Scope, 0, len(scopeMap))
	for scope := range scopeMap {
		scopes = append(scopes, scope)
	}

	return scopes, nil
}

// fetchUsernames lists existing usernames from their index markers.
func fetchUsernames(tr fdb.ReadTransaction) ([]string, error) {
	usernames := []string{}

	iter := tr.GetRange(indexUserSub, fdb.RangeOptions{}).Iterator()
	for iter.Advance() {
		kv := iter.MustGet()

		t, err := indexUserSub.Unpack(kv.Key)
		if err != nil {
			continue
		}

		usernames = append(usernames, t[0].(string))
	}

	return usernames, nil
}
