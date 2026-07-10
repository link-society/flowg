package kv

import (
	"context"
	"io"
)

// Generic Key-Value store, parametrized by transaction types.
type Adapter[QTx QueryTx, MTx MutationTx] interface {
	// Backup streams an incremental snapshot of the database to w, returning the
	// version up to which data was written so it can be passed as since on a
	// subsequent call.
	Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error)

	// Restore loads a snapshot previously produced by Backup from r.
	Restore(ctx context.Context, r io.Reader) error

	// View executes a read-only transaction.
	View(ctx context.Context, fn func(txn QTx) error) error

	// Update executes a read-write transaction.
	Update(ctx context.Context, fn func(txn MTx) error) error
}
