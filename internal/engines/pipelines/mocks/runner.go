package mocks

import (
	"context"

	"io"

	"github.com/stretchr/testify/mock"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
)

// MockRunner is a testify mock implementation of Runner for use in tests.
type MockRunner struct {
	mock.Mock
}

var _ pipelines.Runner = (*MockRunner)(nil)

// NewMockRunner returns a Runner whose calls can be stubbed and asserted.
func NewMockRunner() pipelines.Runner {
	return &MockRunner{}
}

func (m *MockRunner) Run(ctx context.Context, pipelineName string) error {
	args := m.Called(ctx, pipelineName)
	return args.Error(0)
}

func (m *MockRunner) Terminate(ctx context.Context, pipelineName string) error {
	args := m.Called(ctx, pipelineName)
	return args.Error(0)
}

func (m *MockRunner) Process(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error {
	args := m.Called(ctx, pipelineName, entrypoint, record)
	return args.Error(0)
}

func (m *MockRunner) ScrapMetrics(ctx context.Context, pipelineName string, w io.Writer) error {
	args := m.Called(ctx, pipelineName, w)
	return args.Error(0)
}
