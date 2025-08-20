package log

import (
	"context"
	"io"

	"time"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/log/transactions"
	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/utils/ffi/filterdsl"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"
)

type Storage interface {
	proctree.Process
	storage.Streamable

	ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error)
	ListStreamFields(ctx context.Context, stream string) ([]string, error)
	GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error)
	ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error
	DeleteStream(ctx context.Context, stream string) error

	IndexField(ctx context.Context, stream, field string) error
	UnindexField(ctx context.Context, stream, field string) error

	Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error)

	FetchLogs(
		ctx context.Context,
		stream string,
		from, to time.Time,
		filter filterdsl.Filter,
	) ([]models.LogRecord, error)
}

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

type storageImpl struct {
	proctree.Process

	kvStore kvstore.Storage
}

var _ Storage = (*storageImpl)(nil)

func NewStorage(opts ...func(*options)) Storage {
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

	return &storageImpl{
		Process: process,
		kvStore: kvStore,
	}
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.kvStore.Backup(ctx, w, since)
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *storageImpl) ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error) {
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

func (s *storageImpl) ListStreamFields(ctx context.Context, stream string) ([]string, error) {
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

func (s *storageImpl) GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error) {
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

func (s *storageImpl) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.ConfigureStream(txn, stream, config)
		},
	)
}

func (s *storageImpl) DeleteStream(ctx context.Context, stream string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteStream(txn, stream)
		},
	)
}

func (s *storageImpl) IndexField(ctx context.Context, stream, field string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.IndexField(txn, stream, field)
		},
	)
}

func (s *storageImpl) UnindexField(ctx context.Context, stream, field string) error {
	return s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.UnindexField(txn, stream, field)
		},
	)
}

func (s *storageImpl) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error) {
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

func (s *storageImpl) FetchLogs(
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
