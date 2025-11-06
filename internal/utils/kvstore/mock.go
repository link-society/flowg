package kvstore

import (
	"github.com/stretchr/testify/mock"

	"context"
	"io"

	"github.com/dgraph-io/badger/v4"
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

func (m *MockStorage) Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	args := m.Called(ctx, w, since)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockStorage) Restore(ctx context.Context, r io.Reader) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockStorage) View(ctx context.Context, txnFn func(txn *badger.Txn) error) error {
	args := m.Called(ctx, txnFn)
	return args.Error(0)
}

func (m *MockStorage) Update(ctx context.Context, txnFn func(txn *badger.Txn) error) error {
	args := m.Called(ctx, txnFn)
	return args.Error(0)
}
