package mocks

import (
	"github.com/stretchr/testify/mock"

	"context"
	"io"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
)

// MockConfigStorage is a testify-based mock implementation of [ConfigStorage].
type MockConfigStorage struct {
	mock.Mock
}

var _ storage.ConfigStorage = (*MockConfigStorage)(nil)

// NewMockConfigStorage returns a new, unconfigured [MockConfigStorage].
func NewMockConfigStorage() storage.ConfigStorage {
	return &MockConfigStorage{}
}

func (m *MockConfigStorage) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockConfigStorage) Load(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockConfigStorage) ListTransformers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockConfigStorage) ReadTransformer(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockConfigStorage) WriteTransformer(ctx context.Context, name string, content string) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockConfigStorage) DeleteTransformer(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockConfigStorage) ListPipelines(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockConfigStorage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.FlowGraphV2), args.Error(1)
}

func (m *MockConfigStorage) WritePipeline(ctx context.Context, name string, content *models.FlowGraphV2) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockConfigStorage) WriteRawPipeline(ctx context.Context, name string, content string) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockConfigStorage) DeletePipeline(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockConfigStorage) ListForwarders(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockConfigStorage) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.ForwarderV2), args.Error(1)
}

func (m *MockConfigStorage) WriteForwarder(ctx context.Context, name string, content *models.ForwarderV2) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockConfigStorage) DeleteForwarder(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockConfigStorage) HasSystemConfig(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

func (m *MockConfigStorage) ReadSystemConfig(ctx context.Context) (*models.SystemConfiguration, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.SystemConfiguration), args.Error(0)
}

func (m *MockConfigStorage) WriteSystemConfig(ctx context.Context, config *models.SystemConfiguration) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}
