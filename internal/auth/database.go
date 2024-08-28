package auth

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"

	"link-society.com/flowg/internal/logging"
)

type Database struct {
	db *badger.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	opts := badger.
		DefaultOptions(dbPath).
		WithLogger(&logging.BadgerLogger{Channel: "authdb"}).
		WithCompression(options.ZSTD)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
