package mocks

import (
	"github.com/stretchr/testify/mock"

	"context"
	"io"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// MockAuthStorage is a testify-based mock implementation of [AuthStorage].
type MockAuthStorage struct {
	mock.Mock
}

var _ storage.AuthStorage = (*MockAuthStorage)(nil)

// NewMockAuthStorage returns a new, unconfigured [MockAuthStorage].
func NewMockAuthStorage() storage.AuthStorage {
	return &MockAuthStorage{}
}

func (m *MockAuthStorage) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockAuthStorage) Load(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockAuthStorage) ListRoles(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockAuthStorage) FetchRole(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockAuthStorage) SaveRole(ctx context.Context, role models.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockAuthStorage) DeleteRole(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockAuthStorage) ListUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockAuthStorage) FetchUser(ctx context.Context, name string) (*models.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthStorage) ListUserScopes(ctx context.Context, name string) ([]models.Scope, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Scope), args.Error(1)
}

func (m *MockAuthStorage) SaveUser(ctx context.Context, user models.User, password string) error {
	args := m.Called(ctx, user, password)
	return args.Error(0)
}

func (m *MockAuthStorage) PatchUserRoles(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthStorage) DeleteUser(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockAuthStorage) VerifyUserPassword(ctx context.Context, name string, password string) (bool, error) {
	args := m.Called(ctx, name, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthStorage) VerifyUserPermission(ctx context.Context, username string, scope models.Scope) (bool, error) {
	args := m.Called(ctx, username, scope)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthStorage) CreateToken(ctx context.Context, username string) (string, string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthStorage) VerifyToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthStorage) ListTokens(ctx context.Context, username string) ([]string, error) {
	args := m.Called(ctx, username)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthStorage) DeleteToken(ctx context.Context, username string, tokenUUID string) error {
	args := m.Called(ctx, username, tokenUUID)
	return args.Error(0)
}
