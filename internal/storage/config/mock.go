package config

import (
	"github.com/stretchr/testify/mock"

	"context"
	"io"

	"link-society.com/flowg/internal/models"
)

type MockStorage struct {
	mock.Mock
}

var _ Storage = (*MockStorage)(nil)

func NewMockStorage() Storage {
	return &MockStorage{}
}

func (m *MockStorage) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockStorage) Load(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockStorage) ListTransformers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) ReadTransformer(ctx context.Context, name string) (string, error) {
	args := m.Called(ctx, name)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) WriteTransformer(ctx context.Context, name string, content string) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockStorage) DeleteTransformer(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockStorage) ListPipelines(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) ReadPipeline(ctx context.Context, name string) (*models.FlowGraphV2, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.FlowGraphV2), args.Error(1)
}

func (m *MockStorage) WritePipeline(ctx context.Context, name string, content *models.FlowGraphV2) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockStorage) WriteRawPipeline(ctx context.Context, name string, content string) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockStorage) DeletePipeline(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *MockStorage) ListForwarders(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) ReadForwarder(ctx context.Context, name string) (*models.ForwarderV2, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.ForwarderV2), args.Error(1)
}

func (m *MockStorage) WriteForwarder(ctx context.Context, name string, content *models.ForwarderV2) error {
	args := m.Called(ctx, name, content)
	return args.Error(0)
}

func (m *MockStorage) DeleteForwarder(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}
