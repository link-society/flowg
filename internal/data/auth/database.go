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
	opts badger.Options
	db   *badger.DB
}

func NewDatabase(opts DatabaseOpts) *Database {
	var dbDir string
	if !opts.inMemory {
		dbDir = opts.dir
	}

	dbOpts := badger.
		DefaultOptions(dbDir).
		WithLogger(&logging.BadgerLogger{Channel: "authdb"}).
		WithCompression(options.ZSTD).
		WithInMemory(opts.inMemory)

	return &Database{opts: dbOpts, db: nil}
}

func (d *Database) Open() error {
	var err error
	d.db, err = badger.Open(d.opts)
	return err
}

func (d *Database) Close() error {
	if d.db == nil {
		return nil
	}

	return d.db.Close()
}
