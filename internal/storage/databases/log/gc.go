package log

import (
	"log/slog"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/generic/kv"
)

type gcWorker[QTx kv.QueryTx, MTx kv.MutationTx] struct {
	storage    *Storage[QTx, MTx]
	gcInterval time.Duration
}

var _ actor.Worker = (*gcWorker[kv.QueryTx, kv.MutationTx])(nil)

// NewGarbageCollector returns an [actor.Worker] that periodically runs the log
// storage garbage collection, every gcInterval, to enforce each stream's
// retention-size budget. Errors are logged and do not stop the worker.
func NewGarbageCollector[QTx kv.QueryTx, MTx kv.MutationTx](storage *Storage[QTx, MTx], gcInterval time.Duration) actor.Worker {
	return &gcWorker[QTx, MTx]{
		storage:    storage,
		gcInterval: gcInterval,
	}
}

func (w *gcWorker[QTx, MTx]) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.gcInterval):
		go func() {
			if err := w.storage.CollectGarbage(ctx); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to collect garbage",
					slog.String("channel", "logstorage"),
					slog.String("error", err.Error()),
				)
			}
		}()

		return actor.WorkerContinue
	}
}
