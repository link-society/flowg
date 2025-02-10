package kvstore

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/app/logging"

	"link-society.com/flowg/internal/utils/sync"
)

type options struct {
	logChannel string
	dir        string
	inMemory   bool
}

func OptLogChannel(channel string) func(*options) {
	return func(o *options) {
		o.logChannel = channel
	}
}

func OptDirectory(dir string) func(*options) {
	return func(o *options) {
		o.dir = dir
	}
}

func OptInMemory(inMemory bool) func(*options) {
	return func(o *options) {
		o.inMemory = inMemory
	}
}

type Storage struct {
	mbox   actor.Mailbox[message]
	worker *worker
	actor  actor.Actor
	dbOpts badger.Options
}

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		logChannel: "kv",
		dir:        "",
		inMemory:   false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	var dbDir string
	if !options.inMemory {
		dbDir = options.dir
	}

	dbOpts := badger.
		DefaultOptions(dbDir).
		WithLogger(&logging.BadgerLogger{Channel: options.logChannel}).
		WithCompression(badgerOptions.ZSTD).
		WithInMemory(options.inMemory)

	mbox := actor.NewMailbox[message]()
	worker := &worker{
		state: &workerStarting{dbOpts: dbOpts},
		mbox:  mbox,

		startCond: sync.NewCondValue[error](),
		stopCond:  sync.NewCondValue[error](),
	}
	workerA := actor.New(worker)
	actor := actor.Combine(mbox, workerA).WithOptions(actor.OptStopTogether()).Build()

	return &Storage{
		mbox:   mbox,
		worker: worker,
		actor:  actor,
		dbOpts: dbOpts,
	}
}

func (kv *Storage) Start() {
	kv.actor.Start()
}

func (kv *Storage) WaitStarted() error {
	return kv.worker.startCond.Wait()
}

func (kv *Storage) Stop() {
	kv.actor.Stop()
}

func (kv *Storage) WaitStopped() error {
	return kv.worker.stopCond.Wait()
}

func (kv *Storage) View(
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

func (kv *Storage) Update(
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
