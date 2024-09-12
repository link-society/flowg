package auth

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/app/logging"
)

type DatabaseOpts struct {
	dir      string
	inMemory bool
}

func DefaultDatabaseOpts() DatabaseOpts {
	return DatabaseOpts{
		dir:      "./data/auth",
		inMemory: false,
	}
}

func (d DatabaseOpts) WithDir(dir string) DatabaseOpts {
	d.dir = dir
	return d
}

func (d DatabaseOpts) WithInMemory(inMemory bool) DatabaseOpts {
	d.inMemory = inMemory
	return d
}

type Database struct {
	db *badger.DB
}

func NewDatabase(opts DatabaseOpts) (*Database, error) {
	dbOpts := badger.
		DefaultOptions(opts.dir).
		WithLogger(&logging.BadgerLogger{Channel: "authdb"}).
		WithCompression(options.ZSTD).
		WithInMemory(opts.inMemory)

	db, err := badger.Open(dbOpts)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
