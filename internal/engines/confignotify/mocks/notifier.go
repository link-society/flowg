package mocks

import (
	"github.com/stretchr/testify/mock"

	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/confignotify"
)

// MockNotifier is a testify mock implementation of Notifier for use in tests.
type MockNotifier struct {
	mock.Mock
}

var _ confignotify.Notifier = (*MockNotifier)(nil)

// NewMockNotifier returns a Notifier whose calls can be stubbed and asserted.
func NewMockNotifier() confignotify.Notifier {
	return &MockNotifier{}
}

func (m *MockNotifier) Subscribe(ctx context.Context) (actor.MailboxReceiver[confignotify.Event], error) {
	args := m.Called(ctx)
	return args.Get(0).(actor.MailboxReceiver[confignotify.Event]), args.Error(1)
}

func (m *MockNotifier) Notify(ctx context.Context, event confignotify.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
