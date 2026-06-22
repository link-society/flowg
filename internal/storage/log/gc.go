package log

import (
	"log/slog"

	"sync/atomic"
	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/schema"
	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/storage/log/transactions"
)

type gcActor struct {
	actor.Actor
}

type gcWorker struct {
	kvStore    kvstore.Storage
	grace      time.Duration
	gcInterval time.Duration

	running atomic.Bool
}

var _ actor.Worker = (*gcWorker)(nil)

func (w *gcWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.gcInterval):
		if !w.running.CompareAndSwap(false, true) {
			return actor.WorkerContinue
		}

		go func() {
			defer w.running.Store(false)

			if err := w.kvStore.Update(ctx, transactions.CollectGarbage); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to collect garbage",
					slog.String("channel", "logstorage"),
					slog.String("error", err.Error()),
				)
			}

			prefixes := [][]byte{streamConfigPrefix}
			if _, err := schema.CollectGarbage(ctx, w.kvStore, w.grace, prefixes); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to collect tombstones",
					slog.String("channel", "logstorage"),
					slog.String("error", err.Error()),
				)
			}
		}()

		return actor.WorkerContinue
	}
}
