package logstorage

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/app/logging"
)

type Storage struct {
	db *badger.DB
	gc *garbageCollector
}

func NewStorage(dbPath string) (*Storage, error) {
	opts := badger.
		DefaultOptions(dbPath).
		WithLogger(&logging.BadgerLogger{Channel: "logstorage"}).
		WithCompression(options.ZSTD)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	gc := newGarbageCollector(db, 5*time.Minute)
	gc.Start()

	return &Storage{db: db, gc: gc}, nil
}

func (s *Storage) Close() error {
	s.gc.Stop()
	return s.db.Close()
}
