package log

import (
	"context"
	"fmt"
	"log/slog"

	"bytes"
	"encoding/json"
	"io"

	"sync/atomic"
	"time"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/pb"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/changefeed"
	"link-society.com/flowg/internal/storage/log/transactions"
	"link-society.com/flowg/internal/storage/schema"
	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/lww"

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
	Distinct(ctx context.Context, stream string) (map[string][]string, error)

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
	Directory            string
	InMemory             bool
	ReadOnly             bool
	GCInterval           time.Duration
	TombstoneGracePeriod time.Duration
}

const streamKind = "stream"

type storageImpl struct {
	kvStore  kvstore.Storage
	clock    *hlc.Clock
	notifier changefeed.Notifier
}

type deps struct {
	fx.In

	S        kvstore.Storage `name:"storage.log"`
	Clock    *hlc.Clock
	Notifier changefeed.Notifier
}

var _ Storage = (*storageImpl)(nil)

func DefaultOptions() Options {
	return Options{
		Directory:            "",
		InMemory:             false,
		ReadOnly:             false,
		GCInterval:           5 * time.Minute,
		TombstoneGracePeriod: 24 * time.Hour,
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
				kvStore:  d.S,
				clock:    d.Clock,
				notifier: d.Notifier,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					prefixes := [][]byte{[]byte("stream:config:")}
					if err := schema.Migrate(ctx, storage.kvStore, storage.clock, prefixes); err != nil {
						return fmt.Errorf("failed to migrate schema: %w", err)
					}

					return nil
				},
			})

			return storage
		}),
		fxproviders.ProvideActor[*gcActor](func(d deps) *gcActor {
			return &gcActor{
				Actor: actor.New(&gcWorker{
					kvStore:    d.S,
					grace:      opts.TombstoneGracePeriod,
					gcInterval: opts.GCInterval,
				}),
			}
		}),
	)
}

func (s *storageImpl) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.kvStore.Backup(ctx, w, since)
}

func (s *storageImpl) LatestVersion(ctx context.Context) (uint64, error) {
	return s.kvStore.LatestVersion(ctx)
}

func (s *storageImpl) Load(ctx context.Context, r io.Reader) error {
	return s.kvStore.Restore(ctx, r)
}

func (s *storageImpl) Merge(ctx context.Context, r io.Reader) error {
	var changed atomic.Bool
	if err := s.kvStore.Merge(ctx, r, mergeRecord(&changed)); err != nil {
		return err
	}

	if changed.Load() {
		s.emitResync(ctx)
	}
	return nil
}

func (s *storageImpl) ApplyReplicated(ctx context.Context, records []changefeed.Record) error {
	var applied bool
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			for _, record := range records {
				_, ok, err := schema.ApplyRecord(txn, record.Key, record.Value)
				if err != nil {
					return err
				}
				if ok {
					applied = true
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	if applied {
		s.emitResync(ctx)
	}
	return nil
}

var streamConfigPrefix = []byte("stream:config:")

func mergeRecord(changed *atomic.Bool) func(txn *badger.Txn, kv *pb.KV) error {
	return func(txn *badger.Txn, kv *pb.KV) error {
		switch {
		case schema.IsVersionKey(kv.Key):
			return nil

		case bytes.HasPrefix(kv.Key, streamConfigPrefix):
			applied, err := schema.ApplyEnvelope(txn, kv.Key, kv.Value)
			if err != nil {
				return err
			}
			if applied {
				changed.Store(true)
			}
			return nil

		default:
			entry := &badger.Entry{
				Key:       kv.Key,
				Value:     kv.Value,
				ExpiresAt: kv.ExpiresAt,
			}
			if len(kv.UserMeta) > 0 {
				entry.UserMeta = kv.UserMeta[0]
			}
			if err := txn.SetEntry(entry); err != nil {
				return err
			}
			changed.Store(true)
			return nil
		}
	}
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

	ts := s.clock.Now()
	err := s.kvStore.Update(ctx,
		func(txn *badger.Txn) error {
			var err error
			streamConfig, err = transactions.GetOrCreateStreamConfig(txn, stream, ts)
			return err
		},
	)

	if err != nil {
		return models.StreamConfig{}, err
	}

	return streamConfig, nil
}

func (s *storageImpl) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	ts := s.clock.Now()
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.ConfigureStream(txn, stream, config, ts)
		},
	)
	if err != nil {
		return err
	}

	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}
	configVal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}
	record := changefeed.Record{
		Key:   fmt.Appendf(nil, "stream:config:%s", stream),
		Value: lww.Envelope{Timestamp: ts, Payload: configVal}.Marshal(),
	}

	s.emitChange(ctx, stream, changefeed.OpWrite, ts.NodeID, []changefeed.Record{record})
	return nil
}

func (s *storageImpl) DeleteStream(ctx context.Context, stream string) error {
	ts := s.clock.Now()
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.DeleteStream(txn, stream, ts)
		},
	)
	if err != nil {
		return err
	}

	record := changefeed.Record{
		Key:   fmt.Appendf(nil, "stream:config:%s", stream),
		Value: lww.Envelope{Timestamp: ts, Deleted: true}.Marshal(),
	}

	s.emitChange(ctx, stream, changefeed.OpDelete, ts.NodeID, []changefeed.Record{record})
	return nil
}

func (s *storageImpl) emitChange(
	ctx context.Context,
	stream string,
	op changefeed.Operation,
	origin string,
	records []changefeed.Record,
) {
	event := changefeed.ChangeEvent{
		Namespace: changefeed.NamespaceLog,
		Kind:      streamKind,
		Name:      stream,
		Op:        op,
		Origin:    origin,
		Records:   records,
	}
	if err := s.notifier.Notify(ctx, event); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to emit change event",
			"channel", "storage.log",
			"stream", stream,
			"error", err.Error(),
		)
	}
}

func (s *storageImpl) emitResync(ctx context.Context) {
	event := changefeed.ChangeEvent{
		Namespace: changefeed.NamespaceLog,
		Resync:    true,
	}
	if err := s.notifier.Notify(ctx, event); err != nil {
		slog.ErrorContext(
			ctx,
			"failed to emit resync event",
			"channel", "storage.log",
			"error", err.Error(),
		)
	}
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

func (s *storageImpl) Distinct(ctx context.Context, stream string) (map[string][]string, error) {
	var indices map[string][]string

	err := s.kvStore.View(
		ctx,
		func(txn *badger.Txn) error {
			var err error
			indices, err = transactions.Distinct(txn, stream)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return indices, nil
}

func (s *storageImpl) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) ([]byte, error) {
	key := logRecord.NewDbKey(stream)
	ts := s.clock.Now()
	err := s.kvStore.Update(
		ctx,
		func(txn *badger.Txn) error {
			return transactions.Ingest(txn, stream, logRecord, key, ts)
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
