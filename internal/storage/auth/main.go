package auth

import (
	"context"
	"io"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/auth/transactions"
	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/utils/proctree"
)

type Storage interface {
	proctree.Process
	storage.Streamable

	ListRoles(ctx context.Context) ([]models.Role, error)
	FetchRole(ctx context.Context, name string) (*models.Role, error)
	SaveRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, name string) error

	ListUsers(ctx context.Context) ([]models.User, error)
	FetchUser(ctx context.Context, name string) (*models.User, error)
	ListUserScopes(ctx context.Context, name string) ([]models.Scope, error)
	SaveUser(ctx context.Context, user models.User, password string) error
	PatchUserRoles(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, name string) error

	VerifyUserPassword(ctx context.Context, name, password string) (bool, error)
	VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error)

	CreateToken(ctx context.Context, username string) (string, string, error)
	VerifyToken(ctx context.Context, token string) (*models.User, error)
	ListTokens(ctx context.Context, username string) ([]string, error)
	DeleteToken(ctx context.Context, username string, tokenUUID string) error
}

type options struct {
	dir      string
	inMemory bool
	readOnly bool
}

func OptDirectory(dir string) func(*options) {
	return func(o *options) {
		o.dir = dir
	}
}

func OptInMemory(inMemory bool) func(*options) {
	return func(o *options) {
		o.inMemory = inMemory
	}
}

func OptReadOnly(readOnly bool) func(*options) {
	return func(o *options) {
		o.readOnly = readOnly
	}
}

type storageImpl struct {
	proctree.Process

	kvStore *kvstore.Storage
}

var _ Storage = (*storageImpl)(nil)

func NewStorage(opts ...func(*options)) Storage {
	options := options{
		dir:      "",
		inMemory: false,
		readOnly: false,
	}
	for _, opt := range opts {
		opt(&options)
	}

	kvStore := kvstore.NewStorage(
		kvstore.OptLogChannel("authstorage"),
		kvstore.OptDirectory(options.dir),
		kvstore.OptInMemory(options.inMemory),
		kvstore.OptReadOnly(options.readOnly),
	)

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		kvStore,
		proctree.NewProcess(&migratorProcH{kvStore: kvStore}),
	)

	return &storageImpl{
		Process: process,
		kvStore: kvStore,
	}
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
