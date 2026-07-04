package log

import (
	"log/slog"
	"time"

	"github.com/vladopajic/go-actor/actor"

	fdbkvstore "link-society.com/flowg/internal/storage/backends/foundationdb/kvstore"

	"link-society.com/flowg/internal/storage/backends/foundationdb/concrete/log/transactions"
)

type gcWorker struct {
	kvStore    fdbkvstore.Storage
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
