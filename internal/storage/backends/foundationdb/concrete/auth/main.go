package auth

import (
	"context"
	"io"

	"go.uber.org/fx"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	foundationdb_kvstore "link-society.com/flowg/internal/storage/backends/foundationdb/kvstore"

	"link-society.com/flowg/internal/storage/backends/foundationdb/concrete/auth/transactions"
)

// Options configures the foundationdb-backed authentication storage.
type Options struct {
	ConnectionString string
	Prefix           []byte
	InMemory         bool
}

type storageImpl struct {
	kvStore foundationdb_kvstore.Storage
}

type deps struct {
	fx.In

	S foundationdb_kvstore.Storage `name:"storage.auth"`
}

var _ storage.AuthStorage = (*storageImpl)(nil)

// DefaultOptions returns the default [Options] for the authentication storage.
func DefaultOptions() Options {
	return Options{
		ConnectionString: "",
		Prefix:           []byte("flowg/auth"),
		InMemory:         false,
	}
}

// NewStorage returns an fx module providing a foundationdb-backed
// [storage.AuthStorage] configured with the given options.
func NewStorage(opts Options) fx.Option {
	kvOpts := foundationdb_kvstore.DefaultOptions()
	kvOpts.LogChannel = "storage.auth"
	kvOpts.ConnectionString = opts.ConnectionString
	kvOpts.Prefix = opts.Prefix
	kvOpts.InMemory = opts.InMemory

	return fx.Module(
		"storage.auth",
		foundationdb_kvstore.NewStorage(kvOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) storage.AuthStorage {
			// Initialize subspaces from the configured prefix
			root := subspace.FromBytes(opts.Prefix)
			transactions.Init(root)

			impl := &storageImpl{kvStore: d.S}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// Migration not yet needed for FDB; placeholder for future use.
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
		func(tr fdb.ReadTransaction) error {
			var err error
			roles, err = transactions.ListRoles(tr)
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
		func(tr fdb.ReadTransaction) error {
			var err error
			role, err = transactions.FetchRole(tr, name)
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
		func(tr fdb.Transaction) error {
			return transactions.SaveRole(tr, role)
		},
	)
}

func (s *storageImpl) DeleteRole(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(tr fdb.Transaction) error {
			return transactions.DeleteRole(tr, name)
		},
	)
}

func (s *storageImpl) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := s.kvStore.View(
		ctx,
		func(tr fdb.ReadTransaction) error {
			var err error
			users, err = transactions.ListUsers(tr)
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
		func(tr fdb.ReadTransaction) error {
			var err error
			user, err = transactions.FetchUser(tr, name)
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
		func(tr fdb.ReadTransaction) error {
			var err error
			scopes, err = transactions.ListUserScopes(tr, name)
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
		func(tr fdb.Transaction) error {
			return transactions.SaveUser(tr, user, password)
		},
	)
}

func (s *storageImpl) PatchUserRoles(ctx context.Context, user models.User) error {
	return s.kvStore.Update(
		ctx,
		func(tr fdb.Transaction) error {
			return transactions.PatchUserRoles(tr, user)
		},
	)
}

func (s *storageImpl) DeleteUser(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(tr fdb.Transaction) error {
			return transactions.DeleteUser(tr, name)
		},
	)
}

func (s *storageImpl) VerifyUserPassword(ctx context.Context, name, password string) (bool, error) {
	var verified bool

	err := s.kvStore.View(
		ctx,
		func(tr fdb.ReadTransaction) error {
			var err error
			verified, err = transactions.VerifyUserPassword(tr, name, password)
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
		func(tr fdb.ReadTransaction) error {
			var err error
			authorized, err = transactions.VerifyUserPermission(tr, username, scope)
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
		func(tr fdb.Transaction) error {
			var err error
			token, tokenUuid, err = transactions.CreateToken(tr, username)
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
		func(tr fdb.ReadTransaction) error {
			var err error
			user, err = transactions.VerifyToken(tr, token)
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
		func(tr fdb.ReadTransaction) error {
			tokens = transactions.ListTokens(tr, username)
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
		func(tr fdb.Transaction) error {
			return transactions.DeleteToken(tr, username, tokenUUID)
		},
	)
}
