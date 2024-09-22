package lognotify

import (
	"errors"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
)

type subscriberSet map[actor.MailboxSender[LogMessage]]struct{}

type worker struct {
	subscribers map[string]subscriberSet
	subMbox     actor.MailboxReceiver[SubscribeMessage]
	logMbox     actor.MailboxReceiver[LogMessage]
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.subMbox.ReceiveC():
		if !ok {
			go func() {
				msg.ReadyC <- ReadyResponse{errors.New("mailbox closed")}
				close(msg.ReadyC)
			}()
			return actor.WorkerEnd
		}

		if _, exists := w.subscribers[msg.Stream]; !exists {
			w.subscribers[msg.Stream] = make(subscriberSet)
		}

		w.subscribers[msg.Stream][msg.SenderM] = struct{}{}

		go func() {
			<-msg.DoneC
			delete(w.subscribers[msg.Stream], msg.SenderM)
			msg.SenderM.Stop()
		}()

		go func() {
			msg.ReadyC <- ReadyResponse{nil}
			close(msg.ReadyC)
		}()

		return actor.WorkerContinue

	case msg, ok := <-w.logMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if _, exists := w.subscribers[msg.Stream]; !exists {
			w.subscribers[msg.Stream] = make(subscriberSet)
		}

		for sub := range w.subscribers[msg.Stream] {
			if err := sub.Send(ctx, msg); err != nil {
				slog.ErrorContext(
					ctx,
					"failed to send log message",
					"channel", "lognotify",
					"stream", msg.Stream,
					"error", err.Error(),
				)
			}
		}

		return actor.WorkerContinue
	}
}
