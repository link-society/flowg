package auth

import (
	"log/slog"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/schema"
	"link-society.com/flowg/internal/utils/kvstore"
)

type gcActor struct {
	actor.Actor
}

type gcWorker struct {
	kvStore    kvstore.Storage
	grace      time.Duration
	gcInterval time.Duration
}

var _ actor.Worker = (*gcWorker)(nil)

func (w *gcWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.gcInterval):
		go func() {
			if _, err := schema.CollectGarbage(ctx, w.kvStore, w.grace, nil); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to collect tombstones",
					slog.String("channel", "authstorage"),
					slog.String("error", err.Error()),
				)
			}
		}()

		return actor.WorkerContinue
	}
}
