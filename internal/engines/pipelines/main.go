package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type Runner struct {
	mbox actor.MailboxSender[message]

	rootA actor.Actor
}

func NewRunner(
	configStorage *config.Storage,
	logStorage *log.Storage,
	logNotifier *lognotify.LogNotifier,
) *Runner {
	mbox := actor.NewMailbox[message]()
	workerA := actor.New(&worker{
		mbox: mbox,

		configStorage: configStorage,
		logStorage:    logStorage,
		logNotifier:   logNotifier,
	})

	rootA := actor.Combine(mbox, workerA).
		WithOptions(actor.OptStopTogether()).
		Build()

	return &Runner{
		mbox:  mbox,
		rootA: rootA,
	}
}

func (r *Runner) Start() {
	r.rootA.Start()
}

func (r *Runner) Stop() {
	r.rootA.Stop()
}

func (r *Runner) Run(
	ctx context.Context,
	pipelineName string,
	entrypoint string,
	record *models.LogRecord,
) error {
	replyTo := make(chan error)

	err := r.mbox.Send(ctx, message{
		replyTo: replyTo,

		pipelineName: pipelineName,
		entrypoint:   entrypoint,
		record:       record,
	})
	if err != nil {
		return err
	}

	return <-replyTo
}
