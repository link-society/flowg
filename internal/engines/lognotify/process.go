package lognotify

import (
	"errors"
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
)

type subscriberSet map[actor.MailboxSender[LogMessage]]struct{}

type procHandler struct {
	subscribers map[string]subscriberSet
	subMbox     actor.MailboxReceiver[SubscribeMessage]
	logMbox     actor.MailboxReceiver[LogMessage]
}

var _ proctree.ProcessHandler = (*procHandler)(nil)

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case msg, ok := <-h.subMbox.ReceiveC():
		if !ok {
			go func() {
				msg.ReadyC <- ReadyResponse{errors.New("mailbox closed")}
				close(msg.ReadyC)
			}()
			return proctree.Terminate(nil)
		}

		if _, exists := h.subscribers[msg.Stream]; !exists {
			h.subscribers[msg.Stream] = make(subscriberSet)
		}

		h.subscribers[msg.Stream][msg.SenderM] = struct{}{}

		go func() {
			<-msg.DoneC
			delete(h.subscribers[msg.Stream], msg.SenderM)
			msg.SenderM.Stop()
		}()

		go func() {
			msg.ReadyC <- ReadyResponse{nil}
			close(msg.ReadyC)
		}()

		return proctree.Continue()

	case msg, ok := <-h.logMbox.ReceiveC():
		if !ok {
			return proctree.Terminate(nil)
		}

		if _, exists := h.subscribers[msg.Stream]; !exists {
			h.subscribers[msg.Stream] = make(subscriberSet)
		}

		for sub := range h.subscribers[msg.Stream] {
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

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
