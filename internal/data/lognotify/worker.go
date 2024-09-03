package lognotify

import "github.com/vladopajic/go-actor/actor"

type worker struct {
	subscribers map[string]map[chan<- LogMessage]struct{}
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
			w.subscribers[msg.Stream] = make(map[chan<- LogMessage]struct{})
		}

		w.subscribers[msg.Stream][msg.SenderC] = struct{}{}

		go func() {
			<-msg.DoneC
			delete(w.subscribers[msg.Stream], msg.SenderC)
			close(msg.SenderC)
		}()

		return actor.WorkerContinue

	case msg, ok := <-w.logMbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		if _, exists := w.subscribers[msg.Stream]; !exists {
			w.subscribers[msg.Stream] = make(map[chan<- LogMessage]struct{})
		}

		for sub := range w.subscribers[msg.Stream] {
			sub <- msg
		}

		return actor.WorkerContinue
	}
}
