package log

import (
	"github.com/stretchr/testify/mock"

	"context"

	"io"
	"time"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/langs/filterdsl"
)

type MockStorage struct {
	mock.Mock
}

var _ Storage = (*MockStorage)(nil)

func NewMockStorage() Storage {
	return &MockStorage{}
}

func (m *MockStorage) Start() {
	m.Called()
}

func (m *MockStorage) Stop() {
	m.Called()
}

func (m *MockStorage) WaitReady(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Join(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockStorage) Load(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockStorage) ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]models.StreamConfig), args.Error(1)
}

func (m *MockStorage) ListStreamFields(ctx context.Context, stream string) ([]string, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockStorage) GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).(models.StreamConfig), args.Error(1)
}

func (m *MockStorage) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	args := m.Called(ctx, stream, config)
	return args.Error(0)
}

func (m *MockStorage) DeleteStream(ctx context.Context, stream string) error {
	args := m.Called(ctx, stream)
	return args.Error(0)
}

func (m *MockStorage) IndexField(ctx context.Context, stream string, field string) error {
	args := m.Called(ctx, stream, field)
	return args.Error(0)
}

func (m *MockStorage) UnindexField(ctx context.Context, stream string, field string) error {
	args := m.Called(ctx, stream, field)
	return args.Error(0)
}

func (m *MockStorage) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error) {
	args := m.Called(ctx, stream, logRecord)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockStorage) FetchLogs(ctx context.Context, stream string, from time.Time, to time.Time, filter filterdsl.Filter) ([]models.LogRecord, error) {
	args := m.Called(ctx, stream, from, to, filter)
	return args.Get(0).([]models.LogRecord), args.Error(1)
}
