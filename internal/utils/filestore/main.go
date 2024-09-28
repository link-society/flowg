package filestore

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"
)

type options struct {
	dir       string
	inMemory  bool
	extension string
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

func OptExtension(filterExt string) func(*options) {
	return func(o *options) {
		o.extension = filterExt
	}
}

type Storage struct {
	mbox   actor.Mailbox[message]
	worker *worker
	actor  actor.Actor
}

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		dir:       "",
		inMemory:  false,
		extension: "",
	}

	for _, opt := range opts {
		opt(&options)
	}

	mbox := actor.NewMailbox[message]()
	worker := &worker{
		state: &workerStarting{
			baseDir:   options.dir,
			inMemory:  options.inMemory,
			extension: options.extension,
		},
		mbox: mbox,

		startCond: sync.NewCondValue[error](),
		stopCond:  sync.NewCondValue[error](),
	}
	workerA := actor.New(worker)
	actor := actor.Combine(mbox, workerA).WithOptions(actor.OptStopTogether()).Build()

	return &Storage{
		mbox:   mbox,
		worker: worker,
		actor:  actor,
	}
}

func (fs *Storage) Start() {
	fs.actor.Start()
}

func (fs *Storage) WaitStarted() error {
	return fs.worker.startCond.Wait()
}

func (fs *Storage) Stop() {
	fs.actor.Stop()
}

func (fs *Storage) WaitStopped() error {
	return fs.worker.stopCond.Wait()
}

func (fs *Storage) ListFiles(ctx context.Context) ([]string, error) {
	msg := &listItems{
		replyTo: replyTo[[]string]{
			okC:  make(chan []string),
			errC: make(chan error),
		},
	}

	err := fs.mbox.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	return msg.replyTo.Receive()
}

func (fs *Storage) ReadFile(ctx context.Context, key string) ([]byte, error) {
	msg := &readItem{
		replyTo: replyTo[[]byte]{
			okC:  make(chan []byte),
			errC: make(chan error),
		},
		key: key,
	}

	err := fs.mbox.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	return msg.replyTo.Receive()
}

func (fs *Storage) WriteFile(ctx context.Context, key string, content []byte) error {
	msg := &writeItem{
		replyTo: replyTo[bool]{
			okC:  make(chan bool),
			errC: make(chan error),
		},
		key:     key,
		content: content,
	}

	err := fs.mbox.Send(ctx, msg)
	if err != nil {
		return err
	}

	_, err = msg.replyTo.Receive()
	return err
}

func (fs *Storage) DeleteFile(ctx context.Context, key string) error {
	msg := &deleteItem{
		replyTo: replyTo[bool]{
			okC:  make(chan bool),
			errC: make(chan error),
		},
		key: key,
	}

	err := fs.mbox.Send(ctx, msg)
	if err != nil {
		return err
	}

	_, err = msg.replyTo.Receive()
	return err
}
