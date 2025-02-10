package auth

import (
	"context"
	"io"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/storage/auth/transactions"
)

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

type Storage struct {
	kvStore *kvstore.Storage
}

func NewStorage(opts ...func(*options)) *Storage {
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

	return &Storage{kvStore: kvStore}
}

func (s *Storage) Start() {
	s.kvStore.Start()
}

func (s *Storage) WaitStarted() error {
	return s.kvStore.WaitStarted()
}

func (s *Storage) Stop() {
	s.kvStore.Stop()
}

func (s *Storage) WaitStopped() error {
	return s.kvStore.WaitStopped()
}

func (s *Storage) Backup(ctx context.Context, w io.Writer) error {
	return s.kvStore.Backup(ctx, w)
}

func (s *Storage) Restore(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *Storage) ListRoles(ctx context.Context) ([]models.Role, error) {
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

func (s *Storage) FetchRole(ctx context.Context, name string) (*models.Role, error) {
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

func (s *Storage) SaveRole(ctx context.Context, role models.Role) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveRole(txn, role)
		},
	)
}

func (s *Storage) DeleteRole(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteRole(txn, name)
		},
	)
}

func (s *Storage) ListUsers(ctx context.Context) ([]models.User, error) {
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

func (s *Storage) FetchUser(ctx context.Context, name string) (*models.User, error) {
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

func (s *Storage) ListUserScopes(ctx context.Context, name string) ([]models.Scope, error) {
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

func (s *Storage) SaveUser(ctx context.Context, user models.User, password string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.SaveUser(txn, user, password)
		},
	)
}

func (s *Storage) DeleteUser(ctx context.Context, name string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteUser(txn, name)
		},
	)
}

func (s *Storage) VerifyUserPassword(ctx context.Context, name, password string) (bool, error) {
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

func (s *Storage) VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error) {
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

func (s *Storage) CreateToken(ctx context.Context, username string) (string, string, error) {
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

func (s *Storage) VerifyToken(ctx context.Context, token string) (*models.User, error) {
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

func (s *Storage) ListTokens(ctx context.Context, username string) ([]string, error) {
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

func (s *Storage) DeleteToken(ctx context.Context, username string, tokenUUID string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteToken(txn, username, tokenUUID)
		},
	)
}
