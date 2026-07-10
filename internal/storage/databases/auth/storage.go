package auth

import (
	"context"

	"io"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/storage/databases/auth/transactions"
)

// Storage is a backend-agnostic implementation of [storage.AuthStorage]. It
// runs the auth transactions from the transactions subpackage on top of any
// [kv.Adapter], so a single implementation works across every backend.
type Storage[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	adapter kv.Adapter[QTx, MTx]
}

var _ storage.AuthStorage = (*Storage[kv.QueryTx, kv.MutationTx])(nil)

// NewStorage returns a [Storage] that persists auth data through the given
// key-value adapter.
func NewStorage[QTx kv.QueryTx, MTx kv.MutationTx](adapter kv.Adapter[QTx, MTx]) *Storage[QTx, MTx] {
	return &Storage[QTx, MTx]{
		adapter: adapter,
	}
}

// Dump implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) Dump(ctx context.Context, w io.Writer, version uint64) (uint64, error) {
	return s.adapter.Backup(ctx, w, version)
}

// Load implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) Load(ctx context.Context, r io.Reader) error {
	return s.adapter.Restore(ctx, r)
}

// ListRoles implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		roles, err = transactions.ListRoles(txn)
		return err
	})

	return roles, err
}

// FetchRole implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) FetchRole(ctx context.Context, name string) (*models.Role, error) {
	var role *models.Role

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		role, err = transactions.FetchRole(txn, name)
		return err
	})

	return role, err
}

// SaveRole implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) SaveRole(ctx context.Context, role models.Role) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.SaveRole(txn, role)
	})
}

// DeleteRole implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) DeleteRole(ctx context.Context, name string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteRole(txn, name)
	})
}

// ListUsers implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		users, err = transactions.ListUsers(txn)
		return err
	})

	return users, err
}

// FetchUser implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) FetchUser(ctx context.Context, name string) (*models.User, error) {
	var user *models.User

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		user, err = transactions.FetchUser(txn, name)
		return err
	})

	return user, err
}

// ListUserScopes implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) ListUserScopes(ctx context.Context, name string) ([]models.Scope, error) {
	var scopes []models.Scope

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		scopes, err = transactions.ListUserScopes(txn, name)
		return err
	})

	return scopes, err
}

// SaveUser implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) SaveUser(ctx context.Context, user models.User, password string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.SaveUser(txn, user, password)
	})
}

// PatchUserRoles implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) PatchUserRoles(ctx context.Context, user models.User) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.PatchUserRoles(txn, user)
	})
}

// DeleteUser implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) DeleteUser(ctx context.Context, name string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteUser(txn, name)
	})
}

// VerifyUserPassword implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) VerifyUserPassword(ctx context.Context, name string, password string) (bool, error) {
	var verified bool

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		verified, err = transactions.VerifyUserPassword(txn, name, password)
		return err
	})

	return verified, err
}

// VerifyUserPermission implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error) {
	var verified bool

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		verified, err = transactions.VerifyUserPermission(txn, username, scope)
		return err
	})

	return verified, err
}

// CreateToken implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) CreateToken(ctx context.Context, username string) (string, string, error) {
	var token string
	var tokenUUID string

	err := s.adapter.Update(ctx, func(txn MTx) error {
		var err error
		token, tokenUUID, err = transactions.CreateToken(txn, username)
		return err
	})

	return token, tokenUUID, err
}

// VerifyToken implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	var user *models.User

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		user, err = transactions.VerifyToken(txn, token)
		return err
	})

	return user, err
}

// ListTokens implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) ListTokens(ctx context.Context, username string) ([]string, error) {
	var tokens []string

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		tokens = transactions.ListTokens(txn, username)
		return err
	})

	return tokens, err
}

// DeleteToken implements [storage.AuthStorage].
func (s *Storage[QTx, MTx]) DeleteToken(ctx context.Context, username string, tokenUUID string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteToken(txn, username, tokenUUID)
	})
}
