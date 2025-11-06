package kvstore

import (
	"context"
	"fmt"
	"io"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/app/logging"
)

type Storage interface {
	actor.Actor

	Backup(ctx context.Context, w io.Writer, since uint64) (uint64, error)
	Restore(ctx context.Context, r io.Reader) error
	View(ctx context.Context, txnFn func(txn *badger.Txn) error) error
	Update(ctx context.Context, txnFn func(txn *badger.Txn) error) error
}
type Options struct {
	LogChannel string
	Directory  string
	InMemory   bool
	ReadOnly   bool
}
type storageImpl struct {
	actor.Actor

	mbox actor.Mailbox[message]
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		LogChannel: "kv",
		Directory:  "",
		InMemory:   false,
		ReadOnly:   false,
	}
}

func NewStorage(opts Options) fx.Option {
	makeDB := func(lc fx.Lifecycle) (*badger.DB, error) {
		var dbDir string
		if !opts.InMemory {
			dbDir = opts.Directory
		}

		dbOpts := badger.
			DefaultOptions(dbDir).
			WithLogger(&logging.BadgerLogger{Channel: opts.LogChannel}).
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
