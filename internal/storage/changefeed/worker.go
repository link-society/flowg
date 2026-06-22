package changefeed

import (
	"errors"
	"log/slog"

	"sync"

	"github.com/vladopajic/go-actor/actor"
)

type subscriberSet map[actor.MailboxSender[ChangeEvent]]struct{}

type worker struct {
	mu          sync.Mutex
	subscribers subscriberSet
	subMbox     actor.MailboxReceiver[subscribeMessage]
	eventMbox   actor.MailboxReceiver[ChangeEvent]
}

var _ actor.Worker = (*worker)(nil)

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.subMbox.ReceiveC():
		if !ok {
			go func() {
				msg.readyC <- errors.New("mailbox closed")
				close(msg.readyC)
			}()
			return actor.WorkerEnd
		}

		w.mu.Lock()
		w.subscribers[msg.senderM] = struct{}{}
		w.mu.Unlock()

		go func() {
			<-msg.doneC
			w.mu.Lock()
			delete(w.subscribers, msg.senderM)
			w.mu.Unlock()
			msg.senderM.Stop()
		}()

		go func() {
			msg.readyC <- nil
			close(msg.readyC)
		}()

		return actor.WorkerContinue

	case msg, ok := <-w.eventMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		w.mu.Lock()
		subs := make([]actor.MailboxSender[ChangeEvent], 0, len(w.subscribers))
		for sub := range w.subscribers {
			subs = append(subs, sub)
		}
		w.mu.Unlock()

		for _, sub := range subs {
			if err := sub.Send(ctx, msg); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to send change event",
					"channel", "changefeed",
					"namespace", msg.Namespace,
					"error", err.Error(),
				)
			}
		}

		return actor.WorkerContinue
	}
}
