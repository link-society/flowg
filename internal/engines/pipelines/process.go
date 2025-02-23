package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/engines/lognotify"
)

type procHandler struct {
	mbox actor.MailboxReceiver[message]

	configStorage *config.Storage
	logStorage    *log.Storage
	logNotifier   *lognotify.LogNotifier
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case msg, ok := <-h.mbox.ReceiveC():
		if !ok {
			return proctree.Terminate(nil)
		}

		go func() {
			defer close(msg.replyTo)

			pipeline, err := Build(ctx, h.configStorage, msg.pipelineName)
			if err != nil {
				msg.replyTo <- err
				return
			}

			ctx := context.WithValue(ctx, configStorageKey, h.configStorage)
			ctx = context.WithValue(ctx, logStorageKey, h.logStorage)
			ctx = context.WithValue(ctx, logNotifierKey, h.logNotifier)

			err = pipeline.Process(ctx, msg.entrypoint, msg.record)
			msg.replyTo <- err
		}()

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
