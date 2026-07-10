package badger

import (
	"context"
	"fmt"

	"io"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/storage/generic/kv"
)

// Configures how the underlying BadgerDB database is opened.
type AdapterOptions struct {
	// LogChannel names the logging channel and the fx module the store lives in.
	LogChannel string
	// Directory is the on-disk location of the database; ignored when InMemory.
	Directory string
	// InMemory keeps the whole database in memory, useful for tests.
	InMemory bool
	// ReadOnly opens the database without allowing writes.
	ReadOnly bool
}

// BadgerAdapter is a [kv.Adapter] backed by BadgerDB. All database access is
// serialized through an actor mailbox so a single goroutine touches the handle
// at a time.
type BadgerAdapter struct {
	actor.Actor

	mbox actor.Mailbox[message]
}

var _ kv.Adapter[*BadgerTx, *BadgerTx] = (*BadgerAdapter)(nil)

// NewAdapter builds an fx module that provides a [kv.Adapter] backed by
// BadgerDB.
//
// The module wires together the database handle, the actor mailbox and the
// worker that drains it, binding their lifecycles to the fx application. The
// resulting [kv.Adapter] is published under the name given by
// AdapterOptions.LogChannel so several adapters can coexist in the same
// container.
func NewAdapter(opts AdapterOptions) fx.Option {
	makeDB := func(lc fx.Lifecycle) (*badger.DB, error) {
		var dbDir string
		if !opts.InMemory {
			dbDir = opts.Directory
		}

		dbOpts := badger.
			DefaultOptions(dbDir).
			WithLogger(&badgerLogger{Channel: opts.LogChannel}).
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
	) *BadgerAdapter {
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

		adapter := &BadgerAdapter{
			Actor: actor.New(worker),
			mbox:  mbox,
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				adapter.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				adapter.Stop()
				return nil
			},
		})

		return adapter
	}

	module := fmt.Sprintf("kv.adapter.%s", opts.LogChannel)
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

// Backup implements [kv.Adapter.Backup].
func (a *BadgerAdapter) Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	replyTo := make(chan error, 1)
	op := &backupOperation{w: w, since: since}

	err := a.mbox.Send(
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

// Restore implements [kv.Adapter.Restore].
func (a *BadgerAdapter) Restore(ctx context.Context, r io.Reader) error {
	replyTo := make(chan error, 1)

	err := a.mbox.Send(
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

// View implements [kv.Adapter.View].
func (a *BadgerAdapter) View(ctx context.Context, txnFn func(txn *BadgerTx) error) error {
	replyTo := make(chan error, 1)

	err := a.mbox.Send(
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

// Update implements [kv.Adapter.Update].
func (a *BadgerAdapter) Update(ctx context.Context, txnFn func(txn *BadgerTx) error) error {
	replyTo := make(chan error, 1)

	err := a.mbox.Send(
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
