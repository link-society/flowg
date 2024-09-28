package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/engines/lognotify"

	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type worker struct {
	mbox actor.MailboxReceiver[message]

	configStorage *config.Storage
	logStorage    *log.Storage
	logNotifier   *lognotify.LogNotifier
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case msg, ok := <-w.mbox.ReceiveC():
		if !ok {
			return actor.WorkerEnd
		}

		go func() {
			defer close(msg.replyTo)

			pipeline, err := Build(ctx, w.configStorage, msg.pipelineName)
			if err != nil {
				msg.replyTo <- err
				return
			}

			ctx := context.WithValue(ctx, configStorageKey, w.configStorage)
			ctx = context.WithValue(ctx, logStorageKey, w.logStorage)
			ctx = context.WithValue(ctx, logNotifierKey, w.logNotifier)

			err = pipeline.Process(ctx, msg.entrypoint, msg.record)
			msg.replyTo <- err
		}()

		return actor.WorkerContinue
	}
}
