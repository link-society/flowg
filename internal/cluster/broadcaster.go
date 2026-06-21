package cluster

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/changefeed"
)

type broadcaster struct {
	actor.Actor
}

type broadcasterWorker struct {
	localNodeID string

	notifier changefeed.Notifier
	notifyM  actor.MailboxSender[notification]

	eventR actor.MailboxReceiver[changefeed.ChangeEvent]
}

var _ actor.Worker = (*broadcasterWorker)(nil)

func (w *broadcasterWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	if w.eventR == nil {
		eventR, err := w.notifier.Subscribe(ctx)
		if err != nil {
			if ctx.Err() == nil {
				slog.ErrorContext(
					ctx,
					"failed to subscribe to change feed",
					slog.String("channel", "cluster.broadcaster"),
					slog.String("error", err.Error()),
				)
			}

			return actor.WorkerEnd
		}

		w.eventR = eventR
	}

	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case event, ok := <-w.eventR.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if event.Origin != w.localNodeID || len(event.Records) == 0 {
			return actor.WorkerContinue
		}

		msg := &writeNotification{
			Namespace: event.Namespace,
			Records:   event.Records,
		}
		if err := w.notifyM.Send(ctx, msg); err != nil {
			if ctx.Err() == nil {
				slog.ErrorContext(
					ctx,
					"failed to enqueue write notification",
					slog.String("channel", "cluster.broadcaster"),
					slog.String("namespace", event.Namespace),
					slog.String("error", err.Error()),
				)
			}

			return actor.WorkerEnd
		}

		return actor.WorkerContinue
	}
}
