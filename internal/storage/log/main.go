package log

import (
	"context"
	"io"

	"time"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/log/transactions"
	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/utils/langs/filtering"
)

type Storage interface {
	storage.Streamable

	ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error)
	ListStreamFields(ctx context.Context, stream string) ([]string, error)
	GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error)
	ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error
	DeleteStream(ctx context.Context, stream string) error
	StreamUsage(ctx context.Context, stream string) (int64, error)

	IndexField(ctx context.Context, stream, field string) error
	UnindexField(ctx context.Context, stream, field string) error

	Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error)

	FetchLogs(
		ctx context.Context,
		stream string,
		from, to time.Time,
		filter filtering.Filter,
		indexing map[string][]string,
	) ([]models.LogRecord, error)
}

type Options struct {
	Directory  string
	InMemory   bool
	ReadOnly   bool
	GCInterval time.Duration
}

type storageImpl struct {
	kvStore kvstore.Storage
}

type deps struct {
	fx.In

	S kvstore.Storage `name:"storage.log"`
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory:  "",
		InMemory:   false,
		ReadOnly:   false,
		GCInterval: 5 * time.Minute,
	}
}

func NewStorage(opts Options) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "storage.log"
	kvOpts.Directory = opts.Directory
	kvOpts.InMemory = opts.InMemory
	kvOpts.ReadOnly = opts.ReadOnly

	return fx.Module(
		"storage.log",
		kvstore.NewStorage(kvOpts),
		fx.Provide(func(lc fx.Lifecycle, d deps) Storage {
			storage := &storageImpl{
				kvStore: d.S,
			}

			gc := actor.New(&gcWorker{
				kvStore:    d.S,
				gcInterval: opts.GCInterval,
			})

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					gc.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					gc.Stop()
					return nil
				},
			})

			return storage
		}),
	)
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

func (s *storageImpl) StreamUsage(ctx context.Context, stream string) (int64, error) {
	var usage int64

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			usage, err = transactions.EstimateStorage(txn, stream)
			return err
		},
	)

	return usage, err
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
	filter filtering.Filter,
	indexing map[string][]string,
) ([]models.LogRecord, error) {
	var results []models.LogRecord

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			results, err = transactions.FetchLogs(txn, stream, from, to, filter, indexing)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}
