package storage

import (
	"github.com/stretchr/testify/mock"

	"context"

	"io"
	"time"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/langs/filtering"
)

// MockLogStorage is a testify-based mock implementation of [LogStorage].
type MockLogStorage struct {
	mock.Mock
}

var _ LogStorage = (*MockLogStorage)(nil)

// NewMockLogStorage returns a new, unconfigured [MockLogStorage].
func NewMockLogStorage() LogStorage {
	return &MockLogStorage{}
}

func (m *MockLogStorage) Start() {
	m.Called()
}

func (m *MockLogStorage) Stop() {
	m.Called()
}

func (m *MockLogStorage) WaitReady(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogStorage) Join(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockLogStorage) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockLogStorage) Load(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockLogStorage) ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]models.StreamConfig), args.Error(1)
}

func (m *MockLogStorage) ListStreamFields(ctx context.Context, stream string) ([]string, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockLogStorage) GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).(models.StreamConfig), args.Error(1)
}

func (m *MockLogStorage) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	args := m.Called(ctx, stream, config)
	return args.Error(0)
}

func (m *MockLogStorage) DeleteStream(ctx context.Context, stream string) error {
	args := m.Called(ctx, stream)
	return args.Error(0)
}

func (m *MockLogStorage) StreamUsage(ctx context.Context, stream string) (int64, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLogStorage) IndexField(ctx context.Context, stream string, field string) error {
	args := m.Called(ctx, stream, field)
	return args.Error(0)
}

func (m *MockLogStorage) UnindexField(ctx context.Context, stream string, field string) error {
	args := m.Called(ctx, stream, field)
	return args.Error(0)
}

func (m *MockLogStorage) Distinct(ctx context.Context, stream string) (map[string][]string, error) {
	args := m.Called(ctx, stream)
	return args.Get(0).(map[string][]string), args.Error(1)
}

func (m *MockLogStorage) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error) {
	args := m.Called(ctx, stream, logRecord)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockLogStorage) FetchLogs(ctx context.Context, stream string, from time.Time, to time.Time, filter filtering.Filter, indexing map[string][]string) ([]models.LogRecord, error) {
	args := m.Called(ctx, stream, from, to, filter, indexing)
	return args.Get(0).([]models.LogRecord), args.Error(1)
}
