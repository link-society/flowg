package config

import (
	"context"
	"log/slog"

	"maps"

	"crypto/sha256"
	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/confignotify"

	"link-society.com/flowg/internal/storage/backends/foundation"
	"link-society.com/flowg/internal/storage/generic/kv"
)

// watchedItemTypes are the config item prefixes the watcher observes:
// pipeline, transformer and forwarder changes drive the pipeline lifecycle,
// system changes invalidate the cached system configuration.
var watchedItemTypes = []string{"pipeline", "transformer", "forwarder", "system"}

// watchBackoff is how long the watcher waits before re-arming after an error.
const watchBackoff = 5 * time.Second

// snapshot maps item type to item name to a hash of its serialized value.
type snapshot map[string]map[string][sha256.Size]byte

// systemConfigCache is the subset of the config storage the watcher needs to
// drop the cached system configuration on remote changes.
type systemConfigCache interface {
	InvalidateSystemConfigCache()
}

// watchWorker observes the config changelog key and translates mutations into
// confignotify events, so every node of the cluster (including the one that
// made the change) reacts to configuration changes through the same path.
type watchWorker struct {
	adapter     *foundation.FoundationAdapter
	notifier    confignotify.Notifier
	systemCache systemConfigCache

	current snapshot
}

var _ actor.Worker = (*watchWorker)(nil)

// DoWork arms the changelog watch, diffs the config items against the previous
// snapshot, emits the matching events and waits for the watch to fire. Arming
// the watch before scanning guarantees no change is missed in between, and
// watch coalescing is harmless because the diff always compares against the
// stored state. The first scan primes the snapshot silently: the pipeline
// runner builds everything at boot on its own.
func (w *watchWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if ctx.Err() != nil {
		return actor.WorkerEnd
	}

	firedC, err := w.adapter.Watch(ctx, foundation.ChangeLogKey)
	if err != nil {
		w.logError(ctx, "failed to arm config changelog watch", err)
		return w.backoff(ctx)
	}

	next, err := w.scan(ctx)
	if err != nil {
		w.logError(ctx, "failed to scan config items", err)
		return w.backoff(ctx)
	}

	if w.current != nil {
		w.emitDiff(ctx, w.current, next)
	}
	w.current = next

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case err := <-firedC:
		if err != nil {
			if ctx.Err() != nil {
				return actor.WorkerEnd
			}

			w.logError(ctx, "config changelog watch failed", err)
			return w.backoff(ctx)
		}

		return actor.WorkerContinue
	}
}

// scan hashes every watched config item in a single read transaction.
func (w *watchWorker) scan(ctx context.Context) (snapshot, error) {
	next := make(snapshot, len(watchedItemTypes))

	err := w.adapter.View(ctx, func(txn *foundation.FoundationQueryTx) error {
		for _, itemType := range watchedItemTypes {
			items := make(map[string][sha256.Size]byte)

			for pair := range txn.IterPairs(kv.Key{itemType}, kv.KeyRange{}) {
				key := pair.Key()
				items[key[len(key)-1]] = sha256.Sum256(pair.Value())
			}

			next[itemType] = items
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return next, nil
}

// emitDiff broadcasts the difference between two snapshots. A system config
// change drops the local cache; a transformer or forwarder change collapses
// into a single DependenciesChanged event, since the resulting reconcile
// rebuilds every pipeline anyway.
func (w *watchWorker) emitDiff(ctx context.Context, previous snapshot, next snapshot) {
	if !maps.Equal(previous["system"], next["system"]) {
		w.systemCache.InvalidateSystemConfigCache()
	}

	for _, itemType := range []string{"transformer", "forwarder"} {
		if !maps.Equal(previous[itemType], next[itemType]) {
			w.notify(ctx, confignotify.Event{Kind: confignotify.DependenciesChanged})
			return
		}
	}

	for name, hash := range next["pipeline"] {
		if prev, exists := previous["pipeline"][name]; !exists || prev != hash {
			w.notify(ctx, confignotify.Event{
				Kind:     confignotify.PipelineChanged,
				Pipeline: name,
			})
		}
	}

	for name := range previous["pipeline"] {
		if _, exists := next["pipeline"][name]; !exists {
			w.notify(ctx, confignotify.Event{
				Kind:     confignotify.PipelineDeleted,
				Pipeline: name,
			})
		}
	}
}

func (w *watchWorker) notify(ctx context.Context, event confignotify.Event) {
	if err := w.notifier.Notify(ctx, event); err != nil {
		w.logError(ctx, "failed to notify config change", err)
	}
}

// backoff waits before retrying after an error, aborting early on shutdown.
func (w *watchWorker) backoff(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd
	case <-time.After(watchBackoff):
		return actor.WorkerContinue
	}
}

func (w *watchWorker) logError(ctx context.Context, message string, err error) {
	slog.ErrorContext(
		ctx,
		message,
		"channel", "storage.config",
		"error", err.Error(),
	)
}
