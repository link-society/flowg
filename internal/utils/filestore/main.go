package filestore

import (
	"context"

	"io/fs"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
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
	proctree.Process

	mbox actor.MailboxSender[message]
}

var _ proctree.Process = (*Storage)(nil)

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
	handler := &procHandler{
		baseDir:   options.dir,
		inMemory:  options.inMemory,
		extension: options.extension,

		mbox: mbox,
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

func (fs *Storage) StatFile(ctx context.Context, key string) (fs.FileInfo, error) {
	msg := &statItem{
		replyTo: replyTo[statItemResponse]{
			okC:  make(chan statItemResponse),
			errC: make(chan error),
		},
		key: key,
	}

	err := fs.mbox.Send(ctx, msg)
	if err != nil {
		return nil, err
	}

	resp, err := msg.replyTo.Receive()
	if err != nil {
		return nil, err
	}

	return resp.info, nil
}
