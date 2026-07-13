package transactions

import (
	"fmt"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"

	"link-society.com/flowg/internal/utils/hash"
)

// ListUsers loads every user together with their assigned roles.
func ListUsers(txn kv.QueryTx) ([]models.User, error) {
	usernames, err := fetchUsernames(txn)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := []models.User{}

	for _, username := range usernames {
		user, err := FetchUser(txn, username)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user %q: %w", username, err)
		}

		if user != nil {
			users = append(users, *user)
		}
	}

	return users, err
}

// FetchUser loads a single user, returning nil when no "index:user:<name>"
// marker exists. The user's roles are gathered from their "user:<name>:role:*"
// keys.
func FetchUser(txn kv.QueryTx, name string) (*models.User, error) {
	val, err := txn.Get(kv.Key{"index", "user", name})
	if err != nil {
		return nil, fmt.Errorf("failed to read user %q: %w", name, err)
	} else if val == nil {
		return nil, nil
	}

	user := &models.User{Name: name}

	for key := range txn.IterKeys(kv.Key{"user", name, "role"}, kv.KeyRange{}) {
		roleName := key[len(key)-1]
		user.Roles = append(user.Roles, roleName)
	}

	return user, nil
}

// SaveUser stores the user's password hash under "user:<name>:password" and then
// reconciles their role assignments via PatchUserRoles.
func SaveUser(txn kv.MutationTx, user models.User, password string) error {
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password for user %q: %w", user.Name, err)
	}

	key := kv.Key{"user", user.Name, "password"}
	err = txn.Set(key, []byte(passwordHash))
	if err != nil {
		return fmt.Errorf("failed to write password for user %q: %w", user.Name, err)
	}

	return PatchUserRoles(txn, user)
}

// PatchUserRoles writes the "index:user:<name>" marker and converges the stored
// "user:<name>:role:*" keys towards the user's desired role set, adding the
// missing assignments and deleting the obsolete ones (the password is left
// untouched).
func PatchUserRoles(txn kv.MutationTx, user models.User) error {
	key := kv.Key{"index", "user", user.Name}
	err := txn.Set(key, []byte{})
	if err != nil {
		return fmt.Errorf("failed to write index for user %q: %w", user.Name, err)
	}

	obsoleteRoles := map[string]kv.Key{}

	// Collect the roles currently assigned, tentatively flagging any that are no
	// longer part of the desired set as obsolete.
	for key := range txn.IterKeys(kv.Key{"user", user.Name, "role"}, kv.KeyRange{}) {
		roleName := key[len(key)-1]
		if !user.HasRole(roleName) {
			obsoleteRoles[roleName] = key
		}
	}

	// Persist the desired roles: create the ones missing and clear the obsolete
	// flag on the ones already assigned.
	for _, role := range user.Roles {
		if _, exists := obsoleteRoles[role]; !exists {
			key := kv.Key{"user", user.Name, "role", role}
			err := txn.Set(key, []byte{})
			if err != nil {
				return fmt.Errorf("failed to write role %q for user %q: %w", role, user.Name, err)
			}
		} else {
			delete(obsoleteRoles, role)
		}
	}

	// Whatever role keys remain flagged are no longer assigned and get removed.
	for _, key := range obsoleteRoles {
		err := txn.Clear(key)
		if err != nil {
			return fmt.Errorf("failed to clear role %q for user %q: %w", key[len(key)-1], user.Name, err)
		}
	}

	return nil
}

// DeleteUser removes everything tied to a user: all "user:<name>:*" keys (roles
// and password), all of their "pat:<name>:*" tokens (and each token's
// "index:pat:<hash>" reverse-index entry), and the "index:user:<name>" marker.
func DeleteUser(txn kv.MutationTx, name string) error {
	keys := make([]kv.Key, 0)

	for key := range txn.IterKeys(kv.Key{"user", name}, kv.KeyRange{}) {
		keys = append(keys, key)
	}

	for pair := range txn.IterPairs(kv.Key{"pat", name}, kv.KeyRange{}) {
		keys = append(keys, pair.Key())

		if tokenHash := pair.Value(); tokenHash != nil {
			keys = append(keys, kv.Key{"index", "pat", string(tokenHash)})
		}
	}

	for _, key := range keys {
		if err := txn.Clear(key); err != nil {
			return fmt.Errorf("failed to clear key %q for user %q: %w", key, name, err)
		}
	}

	if err := txn.Clear(kv.Key{"index", "user", name}); err != nil {
		return fmt.Errorf("failed to clear index for user %q: %w", name, err)
	}

	return nil
}

// VerifyUserPassword bcrypt-compares a candidate password against the hash
// stored under "user:<name>:password".
func VerifyUserPassword(txn kv.QueryTx, name string, password string) (bool, error) {
	isValid := false

	key := kv.Key{"user", name, "password"}
	val, err := txn.Get(key)
	if err != nil {
		return false, fmt.Errorf("failed to read password for user %q: %w", name, err)
	}

	passwordHash := string(val)
	isValid, err = hash.VerifyPassword(password, passwordHash)
	if err != nil {
		return false, fmt.Errorf("failed to verify password for user %q: %w", name, err)
	}

	return isValid, nil
}

// VerifyUserPermission reports whether a user holds a given scope. It walks the
// user's roles ("user:<name>:role:*") and, for each, the scopes it grants
// ("role:<role>:*"), returning true on a direct match or when the user holds the
// matching write scope for a requested read scope (write implies read).
func VerifyUserPermission(txn kv.QueryTx, name string, scope models.Scope) (bool, error) {
	roles := make([]string, 0)

	for key := range txn.IterKeys(kv.Key{"user", name, "role"}, kv.KeyRange{}) {
		roleName := key[len(key)-1]
		roles = append(roles, roleName)
	}

	for _, roleName := range roles {
		for key := range txn.IterKeys(kv.Key{"role", roleName}, kv.KeyRange{}) {
			scopeName := key[len(key)-1]
			roleScope, err := models.ParseScope(scopeName)
			if err != nil {
				return false, fmt.Errorf("failed to parse scope %q for user %q: %w", scopeName, name, err)
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
func ListUserScopes(txn kv.QueryTx, username string) ([]models.Scope, error) {
	scopeMap := map[models.Scope]struct{}{}
	roles := make([]string, 0)

	for key := range txn.IterKeys(kv.Key{"user", username, "role"}, kv.KeyRange{}) {
		roleName := key[len(key)-1]
		roles = append(roles, roleName)
	}

	for _, roleName := range roles {
		for key := range txn.IterKeys(kv.Key{"role", roleName}, kv.KeyRange{}) {
			scopeName := key[len(key)-1]
			roleScope, err := models.ParseScope(scopeName)
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

// fetchUsernames lists existing usernames from their "index:user:" markers.
func fetchUsernames(txn kv.QueryTx) ([]string, error) {
	usernames := []string{}

	for key := range txn.IterKeys(kv.Key{"index", "user"}, kv.KeyRange{}) {
		username := key[len(key)-1]
		usernames = append(usernames, username)
	}

	return usernames, nil
}
