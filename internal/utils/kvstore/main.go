package kvstore

import (
	"context"
	"io"

	"github.com/vladopajic/go-actor/actor"

	"github.com/dgraph-io/badger/v4"
	badgerOptions "github.com/dgraph-io/badger/v4/options"

	"link-society.com/flowg/internal/app/logging"

	"link-society.com/flowg/internal/utils/proctree"
)

type options struct {
	logChannel string
	dir        string
	inMemory   bool
	readOnly   bool
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

func OptReadOnly(readOnly bool) func(*options) {
	return func(o *options) {
		o.readOnly = readOnly
	}
}

type Storage struct {
	proctree.Process

	mbox actor.Mailbox[message]
}

var _ proctree.Process = (*Storage)(nil)

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		logChannel: "kv",
		dir:        "",
		inMemory:   false,
		readOnly:   false,
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
		WithInMemory(options.inMemory).
		WithReadOnly(options.readOnly)

	mbox := actor.NewMailbox[message]()
	handler := &procHandler{
		dbOpts: dbOpts,
		mbox:   mbox,
	}
	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(mbox),
		proctree.NewProcess(handler),
	)

	return &Storage{
		Process: process,

		mbox: mbox,
	}
}

func (kv *Storage) Backup(
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

func (kv *Storage) Restore(
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
