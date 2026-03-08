package clusterstate

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

var _ Storage = (*MockStorage)(nil)

func NewMockStorage() Storage {
	return &MockStorage{}
}

func (m *MockStorage) FetchLocalState(ctx context.Context, nodeID string, endpoints []string) (*NodeState, error) {
	args := m.Called(ctx, nodeID, endpoints)
	return args.Get(0).(*NodeState), args.Error(1)
}

func (m *MockStorage) UpdateLocalState(ctx context.Context, nodeID string, namespace string, since uint64) error {
	args := m.Called(ctx, nodeID, namespace, since)
	return args.Error(0)
}
