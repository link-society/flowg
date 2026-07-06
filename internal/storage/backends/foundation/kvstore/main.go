package kvstore

import (
	"context"
	"fmt"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// Storage is a concurrency-safe wrapper around a FoundationDB database.
//
// All access goes through an actor mailbox, so transactions issued from
// different goroutines are serialized onto a single worker. This keeps the
// higher-level storage backends free of locking concerns while still allowing
// them to share one underlying database.
type Storage interface {
	actor.Actor

	// View runs txnFn inside a read-only transaction.
	View(ctx context.Context, txnFn func(txn fdb.ReadTransaction) error) error
	// Update runs txnFn inside a read-write transaction.
	Update(ctx context.Context, txnFn func(txn fdb.Transaction) error) error
}

// Options configures how the underlying FoundationDB database is opened.
type Options struct {
	// Tag is a unique identifier for this storage instance, used to name the
	// fx module and the actor mailbox. It is also used to tag the database handle
	// so that multiple storage instances can coexist in the same container.
	// The tag is also used together with the KeySpace to form a unique prefix for
	// all keys stored in this database.
	Tag string

	// ClusterFile is the path to the FoundationDB cluster file, which contains
	// the connection information for the database cluster.
	ClusterFile string
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
		Tag:         "",
		ClusterFile: "",
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
		var (
			db  fdb.Database
			err error
		)

		if opts.ClusterFile == "" {
			db, err = fdb.OpenDefault()
		} else {
			db, err = fdb.OpenDatabase(opts.ClusterFile)
		}
		if err != nil {
			return fdb.Database{}, fmt.Errorf("failed to open FoundationDB database: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
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

	module := fmt.Sprintf("kvstore.%s", opts.Tag)
	tag := func(s string) string { return fmt.Sprintf(`name:"%s.%s"`, module, s) }

	return fx.Module(
		module,
		fx.Provide(
			fx.Annotate(makeDB, fx.ResultTags(tag("db"))),
			fx.Annotate(makeMailbox, fx.ResultTags(tag("mbox"))),
			fx.Annotate(
				makeActor,
				fx.ParamTags("", tag("db"), tag("mbox")),
				fx.ResultTags(fmt.Sprintf(`name:"%s"`, opts.Tag)),
			),
		),
	)
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
