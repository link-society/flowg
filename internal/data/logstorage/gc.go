package logstorage

import (
	"fmt"
	"log/slog"

	"strings"
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
		streams, err := fetchStreamConfigs(txn)
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
				)

				streamSize := int64(0)

				opts := badger.DefaultIteratorOptions
				opts.PrefetchValues = true
				opts.Prefix = []byte(fmt.Sprintf("entry:%s:", stream))
				it := txn.NewIterator(opts)
				defer it.Close()

				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					streamSize += int64(item.EstimatedSize())
				}

				if streamSize > config.RetentionSize {
					slog.DebugContext(
						ctx,
						"Stream too big, purging garbage",
						"channel", "logstorage",
						"stream", stream,
						"size", streamSize,
						"retention_size", config.RetentionSize,
					)

					for it.Rewind(); it.Valid(); it.Next() {
						item := it.Item()
						streamSize -= int64(item.EstimatedSize())

						key := item.KeyCopy(nil)
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

						purgeEntryFromFieldIndex(txn, stream, key)

						if streamSize <= config.RetentionSize {
							break
						}
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

func purgeEntryFromFieldIndex(txn *badger.Txn, stream string, key []byte) error {
	suffix := fmt.Sprintf(":%s", string(key))

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = []byte(fmt.Sprintf("index:%s:field:", stream))
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		indexKey := string(item.Key())

		if strings.HasSuffix(indexKey, suffix) {
			if err := txn.Delete(item.KeyCopy(nil)); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from field index: %w",
					string(key), err,
				)
			}
		}
	}

	return nil
}
