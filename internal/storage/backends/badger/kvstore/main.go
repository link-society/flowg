package kvstore

import (
	"context"
	"fmt"
	"io"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	badgerlog "link-society.com/flowg/internal/storage/backends/badger"
)

// Storage is a concurrency-safe wrapper around a BadgerDB database.
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
	// View runs txnFn inside a read-only transaction.
	View(ctx context.Context, txnFn func(txn *badger.Txn) error) error
	// Update runs txnFn inside a read-write transaction.
	Update(ctx context.Context, txnFn func(txn *badger.Txn) error) error
}

// Options configures how the underlying BadgerDB database is opened.
type Options struct {
	// LogChannel names the logging channel and the fx module the store lives in.
	LogChannel string
	// Directory is the on-disk location of the database; ignored when InMemory.
	Directory string
	// InMemory keeps the whole database in memory, useful for tests.
	InMemory bool
	// ReadOnly opens the database without allowing writes.
	ReadOnly bool
}
type storageImpl struct {
	actor.Actor

	mbox actor.Mailbox[message]
}

var _ Storage = (*storageImpl)(nil)

// DefaultOptions returns the baseline [Options] used when a caller does not
// override them: an on-disk store logging on the "kv" channel.
func DefaultOptions() Options {
	return Options{
		LogChannel: "kv",
		Directory:  "",
		InMemory:   false,
		ReadOnly:   false,
	}
}

// NewStorage builds an fx module that provides a [Storage] backed by BadgerDB.
//
// The module wires together the database handle, the actor mailbox and the
// worker that drains it, binding their lifecycles to the fx application. The
// resulting [Storage] is published under the name given by Options.LogChannel
// so several stores can coexist in the same container.
func NewStorage(opts Options) fx.Option {
	makeDB := func(lc fx.Lifecycle) (*badger.DB, error) {
		var dbDir string
		if !opts.InMemory {
			dbDir = opts.Directory
		}

		dbOpts := badger.
			DefaultOptions(dbDir).
			WithLogger(&badgerlog.BadgerLogger{Channel: opts.LogChannel}).
			WithCompression(badgerOptions.ZSTD).
			WithInMemory(opts.InMemory).
			WithReadOnly(opts.ReadOnly)

		db, err := badger.Open(dbOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return db.Close()
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
		db *badger.DB,
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

		storage := &storageImpl{
			Actor: actor.New(worker),
			mbox:  mbox,
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
	op := &backupOperation{w: w, since: since}

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
			operation: &restoreOperation{r: r},
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
	txnFn func(txn *badger.Txn) error,
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
	txnFn func(txn *badger.Txn) error,
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
