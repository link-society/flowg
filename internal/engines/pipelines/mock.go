package pipelines

import (
	"github.com/stretchr/testify/mock"
	"link-society.com/flowg/internal/models"

	"context"
)

type MockRunner struct {
	mock.Mock
}

var _ Runner = (*MockRunner)(nil)

func NewMockRunner() Runner {
	return &MockRunner{}
}

func (m *MockRunner) Start() {
	m.Called()
}

func (m *MockRunner) Stop() {
	m.Called()
}

func (m *MockRunner) WaitReady(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRunner) Join(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRunner) Run(ctx context.Context, pipelineName string, entrypoint string, record *models.LogRecord) error {
	args := m.Called(ctx, pipelineName, entrypoint, record)
	return args.Error(0)
}
