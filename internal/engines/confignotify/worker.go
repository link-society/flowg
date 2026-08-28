package confignotify

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
)

// worker is the single actor that owns the subscriber registry. All mutations
// happen on the actor goroutine, so the registry needs no locking.
type worker struct {
	subscribers map[actor.Mailbox[Event]]<-chan struct{}
	subMbox     actor.MailboxReceiver[SubscribeMessage]
	eventMbox   actor.MailboxReceiver[Event]
}

var _ actor.Worker = (*worker)(nil)

// DoWork registers one subscriber or broadcasts one event per invocation. Dead
// subscribers (whose done channel fired) are swept on every invocation, on the
// actor goroutine, to avoid mutating the registry concurrently.
func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.subMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		w.sweep()
		w.subscribers[msg.SenderM] = msg.DoneC

		go func() {
			msg.ReadyC <- ReadyResponse{nil}
			close(msg.ReadyC)
		}()

		return actor.WorkerContinue

	case event, ok := <-w.eventMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		w.sweep()

		for sub := range w.subscribers {
			if err := sub.Send(ctx, event); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to deliver config change event",
					"channel", "confignotify",
					"error", err.Error(),
				)
			}
		}

		return actor.WorkerContinue
	}
}

// sweep drops subscribers whose done channel fired and stops their mailboxes.
func (w *worker) sweep() {
	for sub, doneC := range w.subscribers {
		select {
		case <-doneC:
			delete(w.subscribers, sub)
			sub.Stop()
		default:
		}
	}
}
