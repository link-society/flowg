package auth

import (
	"context"
	"fmt"

	"io"

	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/backends/badger/concrete/auth/transactions"
	"link-society.com/flowg/internal/storage/backends/badger/kvstore"
)

// Options configures the badger-backed authentication storage.
type Options struct {
	Directory string
	InMemory  bool
	ReadOnly  bool
}

type storageImpl struct {
	kvStore kvstore.Storage
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"storage.auth"`
}

var _ storage.AuthStorage = (*storageImpl)(nil)

// DefaultOptions returns the default [Options] for the authentication storage.
func DefaultOptions() Options {
	return Options{
		Directory: "",
		InMemory:  false,
		ReadOnly:  false,
	}
}

// NewStorage returns an fx module providing a badger-backed
// [storage.AuthStorage] configured with the given options.
func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "storage.auth"
	kvOpts.Directory = opts.Directory
	kvOpts.InMemory = opts.InMemory
	kvOpts.ReadOnly = opts.ReadOnly

	return fx.Module(
		"storage.auth",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.AuthStorage {
			impl := &storageImpl{kvStore: d.S}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := migrateAlertScopes(ctx, impl.kvStore); err != nil {
						return fmt.Errorf("failed to migrate alerts: %w", err)
					}

					return nil
				},
			})

			return impl
		}),
	)
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.kvStore.Backup(ctx, w, since)
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *storageImpl) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			roles, err = transactions.ListRoles(txn)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *storageImpl) FetchRole(ctx context.Context, name string) (*models.Role, error) {
	var role *models.Role

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			role, err = transactions.FetchRole(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *storageImpl) SaveRole(ctx context.Context, role models.Role) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveRole(txn, role)
		},
	)
}

func (s *storageImpl) DeleteRole(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteRole(txn, name)
		},
	)
}

func (s *storageImpl) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			users, err = transactions.ListUsers(txn)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *storageImpl) FetchUser(ctx context.Context, name string) (*models.User, error) {
	var user *models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			user, err = transactions.FetchUser(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storageImpl) ListUserScopes(ctx context.Context, name string) ([]models.Scope, error) {
	var scopes []models.Scope

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			scopes, err = transactions.ListUserScopes(txn, name)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return scopes, nil
}

func (s *storageImpl) SaveUser(ctx context.Context, user models.User, password string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveUser(txn, user, password)
		},
	)
}

func (s *storageImpl) PatchUserRoles(ctx context.Context, user models.User) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.PatchUserRoles(txn, user)
		},
	)
}

func (s *storageImpl) DeleteUser(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteUser(txn, name)
		},
	)
}

func (s *storageImpl) VerifyUserPassword(ctx context.Context, name, password string) (bool, error) {
	var verified bool

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			verified, err = transactions.VerifyUserPassword(txn, name, password)
			return err
		},
	)
	if err != nil {
		return false, err
	}

	return verified, nil
}

func (s *storageImpl) VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error) {
	var authorized bool

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			authorized, err = transactions.VerifyUserPermission(txn, username, scope)
			return err
		},
	)
	if err != nil {
		return false, err
	}

	return authorized, nil
}

func (s *storageImpl) CreateToken(ctx context.Context, username string) (string, string, error) {
	var token, tokenUuid string

	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			token, tokenUuid, err = transactions.CreateToken(txn, username)
			return err
		},
	)
	if err != nil {
		return "", "", err
	}

	return token, tokenUuid, nil
}

func (s *storageImpl) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	var user *models.User

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			user, err = transactions.VerifyToken(txn, token)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *storageImpl) ListTokens(ctx context.Context, username string) ([]string, error) {
	var tokens []string

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			tokens = transactions.ListTokens(txn, username)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *storageImpl) DeleteToken(ctx context.Context, username string, tokenUUID string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteToken(txn, username, tokenUUID)
		},
	)
}
