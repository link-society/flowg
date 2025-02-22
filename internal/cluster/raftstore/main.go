package raftstore

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/app/logging"
)

type Store struct {
	conn *badger.DB
}

func New(path string) (*Store, error) {
	opts := badger.
		DefaultOptions(path).
		WithLogger(&logging.BadgerLogger{Channel: "raftstore"}).
		WithSyncWrites(true).
		WithCompression(badgerOptions.ZSTD)

	conn, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open raftstore: %w", err)
	}

	return &Store{conn: conn}, nil
}

func (s *Store) Close() error {
	return s.conn.Close()
}
