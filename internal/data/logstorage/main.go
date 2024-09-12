package logstorage

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/app/logging"
)

type StorageOpts struct {
	dir      string
	inMemory bool
}

func DefaultStorageOpts() StorageOpts {
	return StorageOpts{
		dir:      "./data/logs",
		inMemory: false,
	}
}

func (s StorageOpts) WithDir(dir string) StorageOpts {
	s.dir = dir
	return s
}

func (s StorageOpts) WithInMemory(inMemory bool) StorageOpts {
	s.inMemory = inMemory
	return s
}

type Storage struct {
	db      *badger.DB
	gc      *garbageCollector
	indexer *indexer
}

func NewStorage(opts StorageOpts) (*Storage, error) {
	dbOpts := badger.
		DefaultOptions(opts.dir).
		WithLogger(&logging.BadgerLogger{Channel: "logstorage"}).
		WithCompression(options.ZSTD).
		WithInMemory(opts.inMemory)

	db, err := badger.Open(dbOpts)
	if err != nil {
		return nil, err
	}

	gc := newGarbageCollector(db, 5*time.Minute)
	gc.Start()

	indexer := newIndexer(db)
	indexer.Start()

	return &Storage{db: db, gc: gc, indexer: indexer}, nil
}

func (s *Storage) Close() error {
	s.indexer.Stop()
	s.gc.Stop()
	return s.db.Close()
}
