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

// defaultBatchSize bounds how many keys a whole-stream maintenance operation
// (delete, index back-fill, retention eviction) touches per transaction, so a
// single pass stays within a backend's transaction size and time limits.
const defaultBatchSize = 1000

// Storage is a backend-agnostic implementation of [storage.LogStorage]. It runs
// the log transactions from the transactions subpackage on top of any
// [kv.Adapter].
type Storage[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	adapter   kv.Adapter[QTx, MTx]
	batchSize int
}

var _ storage.LogStorage = (*Storage[kv.QueryTx, kv.MutationTx])(nil)

// NewStorage returns a [Storage] that persists log data through the given
// key-value adapter. batchSize bounds how many keys whole-stream operations
// touch per transaction; a value of zero or less uses [defaultBatchSize].
func NewStorage[QTx kv.QueryTx, MTx kv.MutationTx](adapter kv.Adapter[QTx, MTx], batchSize int) *Storage[QTx, MTx] {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}

	return &Storage[QTx, MTx]{
		adapter:   adapter,
		batchSize: batchSize,
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
	var toIndex, toUnindex []string

	err := s.adapter.Update(ctx, func(txn MTx) error {
		var err error
		toIndex, toUnindex, err = transactions.SaveStreamConfig(txn, stream, config)
		return err
	})
	if err != nil {
		return err
	}

	// The index maintenance runs in its own batched transactions, so on a large
	// stream it neither exceeds the backend's limits nor blocks the config write.
	for _, field := range toIndex {
		if err := s.indexField(ctx, stream, field); err != nil {
			return err
		}
	}

	for _, field := range toUnindex {
		if err := s.unindexField(ctx, stream, field); err != nil {
			return err
		}
	}

	return nil
}

// DeleteStream implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) DeleteStream(ctx context.Context, stream string) error {
	for {
		var deleted int

		err := s.adapter.Update(ctx, func(txn MTx) error {
			var err error
			deleted, err = transactions.DeleteStreamDataBatch(txn, stream, s.batchSize)
			return err
		})
		if err != nil {
			return err
		}

		if deleted == 0 {
			break
		}
	}

	// Clear the configuration last so an interrupted deletion can be resumed.
	return s.adapter.Update(ctx, func(txn MTx) error {
		return transactions.DeleteStreamConfig(txn, stream)
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
	return s.indexField(ctx, stream, field)
}

// UnindexField implements [storage.LogStorage].
func (s *Storage[QTx, MTx]) UnindexField(ctx context.Context, stream, field string) error {
	return s.unindexField(ctx, stream, field)
}

// indexField back-fills a field's inverted index over an existing stream, one
// bounded batch of entries per transaction, resuming after the last entry each
// batch processed until the whole stream has been indexed.
func (s *Storage[QTx, MTx]) indexField(ctx context.Context, stream, field string) error {
	var cursor kv.Key

	for {
		var (
			next      kv.Key
			processed int
		)

		err := s.adapter.Update(ctx, func(txn MTx) error {
			var err error
			next, processed, err = transactions.IndexFieldBatch(txn, stream, field, cursor, s.batchSize)
			return err
		})
		if err != nil {
			return err
		}

		if processed == 0 {
			return nil
		}

		cursor = next
	}
}

// unindexField drops a field's inverted index, one bounded batch of keys per
// transaction, until none remain.
func (s *Storage[QTx, MTx]) unindexField(ctx context.Context, stream, field string) error {
	for {
		var deleted int

		err := s.adapter.Update(ctx, func(txn MTx) error {
			var err error
			deleted, err = transactions.UnindexFieldBatch(txn, stream, field, s.batchSize)
			return err
		})
		if err != nil {
			return err
		}

		if deleted == 0 {
			return nil
		}
	}
}

// CollectGarbage enforces every stream's retention-size budget, evicting the
// oldest records of any over-budget stream (and their index references) in
// bounded batches across successive transactions so retention enforcement never
// exceeds a backend's transaction limits.
func (s *Storage[QTx, MTx]) CollectGarbage(ctx context.Context) error {
	var configs map[string]models.StreamConfig

	err := s.adapter.View(ctx, func(txn QTx) error {
		var err error
		configs, err = transactions.FetchStreamConfigs(txn)
		return err
	})
	if err != nil {
		return err
	}

	for stream, config := range configs {
		retentionSize := config.RetentionSize * 1024 * 1024 // MB to bytes
		if retentionSize <= 0 {
			continue
		}

		var streamSize int64
		err := s.adapter.View(ctx, func(txn QTx) error {
			var err error
			streamSize, err = transactions.EstimateStorage(txn, stream)
			return err
		})
		if err != nil {
			return err
		}

		for streamSize > retentionSize {
			var (
				freed   int64
				deleted int
			)

			err := s.adapter.Update(ctx, func(txn MTx) error {
				var err error
				freed, deleted, err = transactions.EvictOldestBatch(
					txn, stream, streamSize-retentionSize, s.batchSize,
				)
				return err
			})
			if err != nil {
				return err
			}

			if deleted == 0 {
				break
			}

			streamSize -= freed
		}
	}

	return nil
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
	key := newDbKey(stream, logRecord)
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
