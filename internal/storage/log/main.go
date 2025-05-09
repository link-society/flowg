package log

import (
	"context"
	"io"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/utils/ffi/filterdsl"
	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/log/transactions"
)

type options struct {
	dir        string
	inMemory   bool
	readOnly   bool
	gcInterval time.Duration
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

func OptGCInterval(interval time.Duration) func(*options) {
	return func(o *options) {
		o.gcInterval = interval
	}
}

type Storage struct {
	proctree.Process

	kvStore *kvstore.Storage
}

var _ proctree.Process = (*Storage)(nil)

func NewStorage(opts ...func(*options)) *Storage {
	options := options{
		dir:        "",
		inMemory:   false,
		readOnly:   false,
		gcInterval: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(&options)
	}

	kvStore := kvstore.NewStorage(
		kvstore.OptLogChannel("logstorage"),
		kvstore.OptDirectory(options.dir),
		kvstore.OptInMemory(options.inMemory),
		kvstore.OptReadOnly(options.readOnly),
	)
	gc := actor.New(&gcWorker{
		kvStore:    kvStore,
		gcInterval: options.gcInterval,
	})

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		kvStore,
		proctree.NewActorProcess(gc),
	)

	return &Storage{
		Process: process,
		kvStore: kvStore,
	}
}

func (s *Storage) Backup(ctx context.Context, w io.Writer) error {
	return s.kvStore.Backup(ctx, w, 0)
}

func (s *Storage) Restore(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *Storage) ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error) {
	var streams map[string]models.StreamConfig

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			streams, err = transactions.FetchStreamConfigs(txn)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return streams, nil
}

func (s *Storage) ListStreamFields(ctx context.Context, stream string) ([]string, error) {
	var fields []string

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			fields = transactions.FetchStreamFields(txn, stream)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return fields, nil
}

func (s *Storage) GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error) {
	var streamConfig models.StreamConfig

	err := s.kvStore.Update(ctx,
		func(txn *badger.Txn) error {
			var err error
			streamConfig, err = transactions.GetOrCreateStreamConfig(txn, stream)
			return err
		},
	)

	if err != nil {
		return models.StreamConfig{}, err
	}

	return streamConfig, nil
}

func (s *Storage) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.ConfigureStream(txn, stream, config)
		},
	)
}

func (s *Storage) DeleteStream(ctx context.Context, stream string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteStream(txn, stream)
		},
	)
}

func (s *Storage) IndexField(ctx context.Context, stream, field string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.IndexField(txn, stream, field)
		},
	)
}

func (s *Storage) UnindexField(ctx context.Context, stream, field string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.UnindexField(txn, stream, field)
		},
	)
}

func (s *Storage) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error) {
	key := logRecord.NewDbKey(stream)
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.Ingest(txn, stream, logRecord, key)
		},
	)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *Storage) FetchLogs(
	ctx context.Context,
	stream string,
	from, to time.Time,
	filter filterdsl.Filter,
) ([]models.LogRecord, error) {
	var results []models.LogRecord

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			results, err = transactions.FetchLogs(txn, stream, from, to, filter)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}
