package cluster

import (
	"context"
	"log/slog"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/changefeed"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type syncTrigger struct {
	actor.Actor
}

type syncTriggerWorker struct {
	logger *slog.Logger

	localNodeID string
	namespaces  map[string]struct{}
	coalesce    time.Duration

	notifier   changefeed.Notifier
	endpoints  *endpointCache
	watermarks *watermarkCache
	requestM   actor.MailboxSender[*syncRequest]

	eventR actor.MailboxReceiver[changefeed.ChangeEvent]
	dirty  map[string]struct{}
	timer  *time.Timer
}

var _ actor.Worker = (*syncTriggerWorker)(nil)

func (w *syncTriggerWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if w.eventR == nil {
		eventR, err := w.notifier.Subscribe(ctx)
		if err != nil {
			if ctx.Err() == nil {
				w.logger.ErrorContext(
					ctx,
					"failed to subscribe to change feed",
					slog.String("channel", "cluster.trigger"),
					slog.String("error", err.Error()),
				)
			}
			return actor.WorkerEnd
		}

		w.eventR = eventR
		w.dirty = make(map[string]struct{})
	}

	var timerC <-chan time.Time
	if w.timer != nil {
		timerC = w.timer.C
	}

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case event, ok := <-w.eventR.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if event.Origin != w.localNodeID {
			return actor.WorkerContinue
		}
		if _, ok := w.namespaces[event.Namespace]; !ok {
			return actor.WorkerContinue
		}

		w.dirty[event.Namespace] = struct{}{}
		// Arm the window once; a continuous write stream still flushes every
		// coalesce interval instead of being starved by repeated resets.
		if w.timer == nil {
			w.timer = time.NewTimer(w.coalesce)
		}
		return actor.WorkerContinue

	case <-timerC:
		w.flush(ctx)
		w.timer = nil
		return actor.WorkerContinue
	}
}

func (w *syncTriggerWorker) flush(ctx context.Context) {
	if len(w.dirty) == 0 {
		return
	}

	namespaces := make([]string, 0, len(w.dirty))
	for ns := range w.dirty {
		namespaces = append(namespaces, ns)
		delete(w.dirty, ns)
	}

	for nodeID, endpoint := range w.endpoints.All() {
		if nodeID == w.localNodeID {
			continue
		}

		lastSync := make([]clusterstate.NamespaceSyncState, 0, len(namespaces))
		for _, ns := range namespaces {
			lastSync = append(lastSync, clusterstate.NamespaceSyncState{
				Namespace: ns,
				Since:     w.watermarks.get(nodeID, ns),
			})
		}

		req := &syncRequest{
			remoteNodeID:   nodeID,
			remoteEndpoint: endpoint,
			lastSync:       lastSync,
		}
		if err := w.requestM.Send(ctx, req); err != nil {
			if ctx.Err() == nil {
				w.logger.ErrorContext(
					ctx,
					"failed to enqueue triggered sync request",
					slog.String("channel", "cluster.trigger"),
					slog.String("cluster.remote.node", nodeID),
					slog.String("error", err.Error()),
				)
			}
		}
	}
}
