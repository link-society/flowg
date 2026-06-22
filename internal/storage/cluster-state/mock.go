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

func (m *MockStorage) GetLiveness(ctx context.Context, namespace string) (int64, error) {
	args := m.Called(ctx, namespace)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) SetLiveness(ctx context.Context, namespace string, unixNano int64) error {
	args := m.Called(ctx, namespace, unixNano)
	return args.Error(0)
}

func (m *MockStorage) ResetLocalState(ctx context.Context, namespace string) error {
	args := m.Called(ctx, namespace)
	return args.Error(0)
}
