package log

import (
	"context"
	"io"

	"time"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/storage/databases/log/transactions"

	"link-society.com/flowg/internal/utils/langs/filtering"
)

// Storage is a backend-agnostic implementation of [storage.LogStorage]. It runs
// the log transactions from the transactions subpackage on top of any
// [kv.Adapter].
type Storage[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	adapter kv.Adapter[QTx, MTx]
}

var _ storage.LogStorage = (*Storage[kv.QueryTx, kv.MutationTx])(nil)

// NewStorage returns a [Storage] that persists log data through the given
// key-value adapter.
func NewStorage[QTx kv.QueryTx, MTx kv.MutationTx](adapter kv.Adapter[QTx, MTx]) *Storage[QTx, MTx] {
	return &Storage[QTx, MTx]{
		adapter: adapter,
	}
}

// Dump implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) Dump(ctx context.Context, w io.Writer, since uint64) (uint64, error) {
	return s.adapter.Backup(ctx, w, since)
}

// Load implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) Load(ctx context.Context, r io.Reader) error {
	return s.adapter.Restore(ctx, r)
}

// ListStreamConfigs implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) ListStreamConfigs(ctx context.Context) (map[string]models.StreamConfig, error) {
	var streams map[string]models.StreamConfig

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		streams, err = transactions.FetchStreamConfigs(txn)
		return err
	})
	if err != nil {
		return nil, err
	}

	return streams, nil
}

// ListStreamFields implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) ListStreamFields(ctx context.Context, stream string) ([]string, error) {
	var fields []string

	err := s.adapter.View(ctx, func(txn QTx) error {
		fields = transactions.FetchStreamFields(txn, stream)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fields, nil
}

// GetOrCreateStreamConfig implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) GetOrCreateStreamConfig(ctx context.Context, stream string) (models.StreamConfig, error) {
	var streamConfig models.StreamConfig

	err := s.adapter.Update(ctx, func(txn MTx) error {
		var err error
		streamConfig, err = transactions.GetOrCreateStreamConfig(txn, stream)
		return err
	})

	if err != nil {
		return models.StreamConfig{}, err
	}

	return streamConfig, nil
}

// ConfigureStream implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) ConfigureStream(ctx context.Context, stream string, config models.StreamConfig) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.ConfigureStream(txn, stream, config)
	})
}

// DeleteStream implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) DeleteStream(ctx context.Context, stream string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteStream(txn, stream)
	})
}

// StreamUsage implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) StreamUsage(ctx context.Context, stream string) (int64, error) {
	var usage int64

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		usage, err = transactions.EstimateStorage(txn, stream)
		return err
	})

	return usage, err
}

// IndexField implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) IndexField(ctx context.Context, stream, field string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.IndexField(txn, stream, field)
	})
}

// UnindexField implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) UnindexField(ctx context.Context, stream, field string) error {
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.UnindexField(txn, stream, field)
	})
}

// Distinct implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) Distinct(ctx context.Context, stream string) (map[string][]string, error) {
	var indices map[string][]string

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		indices, err = transactions.Distinct(txn, stream)
		return err
	})
	if err != nil {
		return nil, err
	}

	return indices, nil
}

// Ingest implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) Ingest(ctx context.Context, stream string, logRecord *models.LogRecord) (kv.Key, error) {
	key := logRecord.NewDbKey(stream)
	err := s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.Ingest(txn, stream, logRecord, key)
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

// FetchLogs implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) FetchLogs(
	ctx context.Context,
	stream string,
	from, to time.Time,
	filter filtering.Filter,
	indexing map[string][]string,
) ([]models.LogRecord, error) {
	var results []models.LogRecord

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		results, err = transactions.FetchLogs(txn, stream, from, to, filter, indexing)
		return err
	})
	if err != nil {
		return nil, err
	}

	return results, nil
}
