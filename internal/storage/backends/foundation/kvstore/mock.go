package kvstore

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
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

func (m *MockStorage) View(ctx context.Context, txnFn func(txn fdb.ReadTransaction) error) error {
	args := m.Called(ctx, txnFn)
	return args.Error(0)
}

func (m *MockStorage) Update(ctx context.Context, txnFn func(txn fdb.Transaction) error) error {
	args := m.Called(ctx, txnFn)
	return args.Error(0)
}
