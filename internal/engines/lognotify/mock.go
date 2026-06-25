package lognotify

import (
	"github.com/stretchr/testify/mock"

	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"
)

// MockNotifier is a testify mock implementation of LogNotifier for use in tests.
type MockNotifier struct {
	mock.Mock
}

var _ LogNotifier = (*MockNotifier)(nil)

// NewMockNotifier returns a LogNotifier whose calls can be stubbed and asserted.
func NewMockNotifier() LogNotifier {
	return &MockNotifier{}
}

func (m *MockNotifier) Subscribe(ctx context.Context, stream string) (actor.MailboxReceiver[LogMessage], error) {
	args := m.Called(ctx, stream)
	return args.Get(0).(actor.MailboxReceiver[LogMessage]), args.Error(1)
}

func (m *MockNotifier) Notify(ctx context.Context, stream string, logKey string, logRecord models.LogRecord) error {
	args := m.Called(ctx, stream, logKey, logRecord)
	return args.Error(0)
}
