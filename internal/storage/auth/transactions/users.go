package transactions

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/utils/auth/hash"
	"link-society.com/flowg/internal/utils/hlc"
)

func ListUsers(txn *badger.Txn) ([]models.User, error) {
	usernames, err := fetchUsernames(txn)
	if err != nil {
		return nil, err
	}

	users := []models.User{}

	for _, username := range usernames {
		user, err := FetchUser(txn, username)
		if err != nil {
			return nil, err
		}

		if user != nil {
			users = append(users, *user)
		}
	}

	return users, err
}

func FetchUser(txn *badger.Txn, name string) (*models.User, error) {
	_, found, err := getItem(txn, []byte(fmt.Sprintf("index:user:%s", name)))
	if err != nil {
		return nil, fmt.Errorf("failed to get index of user '%s': %w", name, err)
	}
	if !found {
		return nil, nil
	}

	user := &models.User{Name: name}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", name))
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
		roleName := string(item.Key()[len(name)+11:])
		user.Roles = append(user.Roles, roleName)
	}

	return user, nil
}

func SaveUser(txn *badger.Txn, user models.User, password string, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	key := []byte(fmt.Sprintf("user:%s:password", user.Name))
	if err := setItem(txn, key, []byte(passwordHash), ts, sink); err != nil {
		return fmt.Errorf("failed to save password of user '%s': %w", user.Name, err)
	}

	return PatchUserRoles(txn, user, ts, sink)
}

func PatchUserRoles(txn *badger.Txn, user models.User, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	indexKey := []byte(fmt.Sprintf("index:user:%s", user.Name))
	if err := setItem(txn, indexKey, []byte{}, ts, sink); err != nil {
		return fmt.Errorf("failed to save index of user '%s': %w", user.Name, err)
	}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", user.Name))
	it := txn.NewIterator(opts)

	obsoleteRoles := map[string][]byte{}

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
		roleName := string(key[len(user.Name)+11:])
		if !user.HasRole(roleName) {
			obsoleteRoles[roleName] = key
		}
	}
	it.Close()

	for _, role := range user.Roles {
		if _, exists := obsoleteRoles[role]; !exists {
			key := []byte(fmt.Sprintf("user:%s:role:%s", user.Name, role))
			if err := setItem(txn, key, []byte{}, ts, sink); err != nil {
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
		if err := deleteItem(txn, key, ts, sink); err != nil {
			return fmt.Errorf(
				"failed to delete role '%s' from user '%s': %w",
				role, user.Name, err,
			)
		}
	}

	return nil
}

func DeleteUser(txn *badger.Txn, name string, ts hlc.Timestamp, sink *[]changefeed.Record) error {
	keys := make([][]byte, 0)

	collect := func(prefix string) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			_, live, err := liveValue(item)
			if err != nil {
				return err
			}
			if live {
				keys = append(keys, item.KeyCopy(nil))
			}
		}
		return nil
	}

	if err := collect(fmt.Sprintf("user:%s:", name)); err != nil {
		return err
	}
	if err := collect(fmt.Sprintf("pat:%s:", name)); err != nil {
		return err
	}

	for _, key := range keys {
		if err := deleteItem(txn, key, ts, sink); err != nil {
			return fmt.Errorf(
				"failed to delete key '%s' of user '%s': %w",
				name, key, err,
			)
		}
	}

	indexKey := []byte(fmt.Sprintf("index:user:%s", name))
	if err := deleteItem(txn, indexKey, ts, sink); err != nil {
		return fmt.Errorf(
			"failed to delete index of user '%s': %w",
			name, err,
		)
	}

	return nil
}

func VerifyUserPassword(txn *badger.Txn, name, password string) (bool, error) {
	isValid := false

	key := []byte(fmt.Sprintf("user:%s:password", name))
	payload, found, err := getItem(txn, key)
	if err != nil {
		return false, fmt.Errorf("failed to get password of user '%s': %w", name, err)
	}
	if !found {
		return false, fmt.Errorf("failed to get password of user '%s': %w", name, badger.ErrKeyNotFound)
	}

	isValid, err = hash.VerifyPassword(password, string(payload))
	if err != nil {
		return false, fmt.Errorf("failed to verify password of user '%s': %w", name, err)
	}

	return isValid, nil
}

func VerifyUserPermission(txn *badger.Txn, name string, scope models.Scope) (bool, error) {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", name))
	it := txn.NewIterator(opts)
	defer it.Close()

	roles := make([]string, 0)

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			return false, err
		}
		if !live {
			continue
		}
		roleName := string(item.Key()[len(name)+11:])
		roles = append(roles, roleName)
	}

	for _, roleName := range roles {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(fmt.Sprintf("role:%s:", roleName))
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			_, live, err := liveValue(item)
			if err != nil {
				return false, err
			}
			if !live {
				continue
			}
			scopeName := string(item.Key()[len(roleName)+6:])
			roleScope, err := models.ParseScope(scopeName)
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

func ListUserScopes(txn *badger.Txn, username string) ([]models.Scope, error) {
	scopeMap := map[models.Scope]struct{}{}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte(fmt.Sprintf("user:%s:role:", username))
	it := txn.NewIterator(opts)
	defer it.Close()

	roles := make([]string, 0)

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		_, live, err := liveValue(item)
		if err != nil {
			return nil, err
		}
		if !live {
			continue
		}
		roleName := string(item.Key()[len(username)+11:])
		roles = append(roles, roleName)
	}

	for _, roleName := range roles {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(fmt.Sprintf("role:%s:", roleName))
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
			scopeName := string(item.Key()[len(roleName)+6:])
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

func fetchUsernames(txn *badger.Txn) ([]string, error) {
	usernames := []string{}

	opts := badger.DefaultIteratorOptions
	opts.Prefix = []byte("index:user:")
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
		usernames = append(usernames, string(item.Key()[11:]))
	}

	return usernames, nil
}
