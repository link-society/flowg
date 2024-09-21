package lognotify

import "github.com/vladopajic/go-actor/actor"

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

		return actor.WorkerContinue

	case msg, ok := <-w.logMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if _, exists := w.subscribers[msg.Stream]; !exists {
			w.subscribers[msg.Stream] = make(subscriberSet)
		}

		for sub := range w.subscribers[msg.Stream] {
			sub.Send(ctx, msg)
		}

		return actor.WorkerContinue
	}
}
