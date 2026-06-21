package changefeed

import (
	"errors"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
)

type subscriberSet map[actor.MailboxSender[ChangeEvent]]struct{}

type worker struct {
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

		w.subscribers[msg.senderM] = struct{}{}

		go func() {
			<-msg.doneC
			delete(w.subscribers, msg.senderM)
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

		for sub := range w.subscribers {
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
