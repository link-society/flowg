package logstorage

import (
	"log/slog"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/vladopajic/go-actor/actor"
)

type garbageCollector struct {
	actor actor.Actor
}

type garbageCollectorWorker struct {
	db                 *badger.DB
	collectionInterval time.Duration
}

func newGarbageCollector(
	db *badger.DB,
	collectionInterval time.Duration,
) *garbageCollector {
	return &garbageCollector{
		actor: actor.New(&garbageCollectorWorker{
			db:                 db,
			collectionInterval: collectionInterval,
		}),
	}
}

func (gc *garbageCollector) Start() {
	gc.actor.Start()
}

func (gc *garbageCollector) Stop() {
	gc.actor.Stop()
}

func (w *garbageCollectorWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case <-time.After(w.collectionInterval):
		w.CollectGarbage(ctx)
		return actor.WorkerContinue
	}
}

func (w *garbageCollectorWorker) CollectGarbage(ctx actor.Context) {
	err := w.db.Update(func(txn *badger.Txn) error {

		return nil
	})

	if err != nil {
		slog.ErrorContext(
			ctx,
			"Failed to collect garbage",
			"channel", "logstorage",
			"error", err.Error(),
		)
	}
}
