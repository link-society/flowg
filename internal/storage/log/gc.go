package log

import (
	"log/slog"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/kvstore"

	"link-society.com/flowg/internal/storage/log/transactions"
)

type gcWorker struct {
	kvStore    *kvstore.Storage
	gcInterval time.Duration
}

var _ actor.Worker = (*gcWorker)(nil)

func (w *gcWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.gcInterval):
		go func() {
			if err := w.kvStore.Update(ctx, transactions.CollectGarbage); err != nil {
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
