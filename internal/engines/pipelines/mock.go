package pipelines

import (
	"github.com/stretchr/testify/mock"
	"link-society.com/flowg/internal/models"

	"context"
)

// MockRunner is a testify mock implementation of Runner for use in tests.
type MockRunner struct {
	mock.Mock
}

var _ Runner = (*MockRunner)(nil)

// NewMockRunner returns a Runner whose calls can be stubbed and asserted.
func NewMockRunner() Runner {
	return &MockRunner{}
}

func (m *MockRunner) Run(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error {
	args := m.Called(ctx, pipelineName, entrypoint, record)
	return args.Error(0)
}

func (m *MockRunner) InvalidateCachedBuild(ctx context.Context, pipelineName string) error {
	args := m.Called(ctx, pipelineName)
	return args.Error(0)
}

func (m *MockRunner) InvalidateAllCachedBuilds(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
