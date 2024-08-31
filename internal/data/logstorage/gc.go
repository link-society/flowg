package logstorage

import (
	"fmt"
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
		streams, err := fetchStreamconfigs(txn)
		if err != nil {
			return err
		}

		for stream, config := range streams {
			if config.RetentionSize > 0 {
				slog.DebugContext(
					ctx,
					"Collecting garbage",
					"channel", "logstorage",
					"stream", stream,
					"retention_size", config.RetentionSize,
				)

				streamSize := int64(0)
				garbageKeys := [][]byte{}

				opts := badger.DefaultIteratorOptions
				opts.PrefetchValues = false
				opts.Prefix = []byte(fmt.Sprintf("entry:%s:", stream))
				opts.Reverse = true
				it := txn.NewIterator(opts)
				defer it.Close()

				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					streamSize += item.EstimatedSize()

					if streamSize > config.RetentionSize {
						garbageKeys = append(garbageKeys, item.KeyCopy(nil))
					}
				}

				for _, key := range garbageKeys {
					slog.DebugContext(
						ctx,
						"Purging key from BadgerDB",
						"channel", "logstorage",
						"stream", stream,
						"key", key,
					)

					if err := txn.Delete(key); err != nil {
						return fmt.Errorf(
							"could not delete key '%s' from stream '%s': %w",
							key, stream, err,
						)
					}
				}
			}
		}

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
