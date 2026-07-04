package kvstore

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
)

// Storage is a concurrency-safe wrapper around a FoundationDB database.
//
// All access goes through an actor mailbox, so transactions issued from
// different goroutines are serialized onto a single worker. This keeps the
// higher-level storage backends free of locking concerns while still allowing
// them to share one underlying database.
type Storage interface {
	actor.Actor

	// Backup streams an incremental snapshot of the database to w, returning the
	// version up to which data was written so it can be passed as since on a
	// subsequent call.
	Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error)
	// Restore loads a snapshot previously produced by Backup from r.
	Restore(ctx context.Context, r io.Reader) error
	// View runs txnFn inside a read-only FDB transaction.
	View(ctx context.Context, txnFn func(txn fdb.ReadTransaction) error) error
	// Update runs txnFn inside a read-write FDB transaction with automatic retry
	// on conflict.
	Update(ctx context.Context, txnFn func(txn fdb.Transaction) error) error
}

// Options configures how the FoundationDB database is opened.
type Options struct {
	// LogChannel names the logging channel and the fx module the store lives in.
	LogChannel string
	// ConnectionString is an optional FDB connection string.
	// When empty, the default cluster file is used.
	ConnectionString string
	// Prefix defines the subspace prefix for backup/restore operations.
	// When nil or empty, the entire database is used.
	Prefix []byte
	// InMemory is a compatibility option.
	// FoundationDB does not support in-memory mode; this flag is ignored.
	InMemory bool
}

type storageImpl struct {
	actor.Actor

	mbox   actor.Mailbox[message]
	db     fdb.Database
	prefix []byte
}

var _ Storage = (*storageImpl)(nil)

var fdbAPIVersionOnce sync.Once

// DefaultOptions returns the baseline [Options] used when a caller does not
// override them.
func DefaultOptions() Options {
	return Options{
		LogChannel:       "kv",
		ConnectionString: "",
		Prefix:           nil,
		InMemory:         false,
	}
}

// NewStorage builds an fx module that provides a [Storage] backed by FoundationDB.
//
// The module wires together the database handle, the actor mailbox and the
// worker that drains it, binding their lifecycles to the fx application. The
// resulting [Storage] is published under the name given by Options.LogChannel
// so several stores can coexist in the same container.
func NewStorage(opts Options) fx.Option {
	makeDB := func(lc fx.Lifecycle) (fdb.Database, error) {
		fdbAPIVersionOnce.Do(func() {
			fdb.MustAPIVersion(730)
		})

		var db fdb.Database
		var err error

		if opts.ConnectionString != "" {
			db, err = fdb.OpenWithConnectionString(opts.ConnectionString)
		} else {
			db, err = fdb.OpenDefault()
		}
		if err != nil {
			return fdb.Database{}, fmt.Errorf("failed to open foundationdb: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				slog.Debug("closing foundationdb database", "channel", opts.LogChannel)
				db.Close()
				return nil
			},
		})

		return db, nil
	}

	makeMailbox := func(lc fx.Lifecycle) actor.Mailbox[message] {
		mbox := actor.NewMailbox[message]()

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				mbox.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				mbox.Stop()
				return nil
			},
		})

		return mbox
	}

	makeActor := func(
		lc fx.Lifecycle,
		db fdb.Database,
		mbox actor.Mailbox[message],
	) Storage {
		worker := actor.NewWorker(func(ctx actor.Context) actor.WorkerStatus {
			select {
			case <-ctx.Done():
				return actor.WorkerEnd

			case msg, ok := <-mbox.ReceiveC():
				if !ok {
					return actor.WorkerEnd
				}

				go func() {
					err := msg.operation.Handle(db)
					msg.replyTo <- err
					close(msg.replyTo)
				}()

				return actor.WorkerContinue
			}
		})

		var prefix []byte
		if opts.Prefix != nil {
			prefix = make([]byte, len(opts.Prefix))
			copy(prefix, opts.Prefix)
		}

		storage := &storageImpl{
			Actor:  actor.New(worker),
			mbox:   mbox,
			db:     db,
			prefix: prefix,
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				storage.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				storage.Stop()
				return nil
			},
		})

		return storage
	}

	module := fmt.Sprintf("kvstore.%s", opts.LogChannel)
	tag := func(s string) string { return fmt.Sprintf(`name:"%s.%s"`, module, s) }

	return fx.Module(
		module,
		fx.Provide(
			fx.Annotate(makeDB, fx.ResultTags(tag("db"))),
			fx.Annotate(makeMailbox, fx.ResultTags(tag("mbox"))),
			fx.Annotate(
				makeActor,
				fx.ParamTags("", tag("db"), tag("mbox")),
				fx.ResultTags(fmt.Sprintf(`name:"%s"`, opts.LogChannel)),
			),
		),
	)
}

// Backup implements [Storage.Backup].
func (kv *storageImpl) Backup(
	ctx context.Context,
	w io.Writer,
	since uint64,
) (uint64, error) {
	replyTo := make(chan error, 1)
	op := &backupOperation{w: w, since: since, prefix: kv.prefix}

	err := kv.mbox.Send(
		ctx,
		message{
			replyTo:   replyTo,
			operation: op,
		},
	)
	if err != nil {
		return 0, err
	}

	err = <-replyTo
	if err != nil {
		return 0, err
	}

	return op.since, nil
}

// Restore implements [Storage.Restore].
func (kv *storageImpl) Restore(
	ctx context.Context,
	r io.Reader,
) error {
	replyTo := make(chan error, 1)

	err := kv.mbox.Send(
		ctx,
		message{
			replyTo:   replyTo,
			operation: &restoreOperation{r: r, prefix: kv.prefix},
		},
	)
	if err != nil {
		return err
	}

	return <-replyTo
}

// View implements [Storage.View].
func (kv *storageImpl) View(
	ctx context.Context,
	txnFn func(txn fdb.ReadTransaction) error,
) error {
	replyTo := make(chan error, 1)

	err := kv.mbox.Send(
		ctx,
		message{
			replyTo:   replyTo,
			operation: &viewOperation{txnFn: txnFn},
		},
	)
	if err != nil {
		return err
	}

	return <-replyTo
}

// Update implements [Storage.Update].
func (kv *storageImpl) Update(
	ctx context.Context,
	txnFn func(txn fdb.Transaction) error,
) error {
	replyTo := make(chan error, 1)

	err := kv.mbox.Send(
		ctx,
		message{
			replyTo:   replyTo,
			operation: &updateOperation{txnFn: txnFn},
		},
	)
	if err != nil {
		return err
	}

	return <-replyTo
}
