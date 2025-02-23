package pipelines

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
)

type Runner struct {
	mbox    actor.MailboxSender[message]
	process proctree.Process
}

func NewRunner(
	configStorage *config.Storage,
	logStorage *log.Storage,
	logNotifier *lognotify.LogNotifier,
) *Runner {
	mbox := actor.NewMailbox[message]()
	handler := &procHandler{
		mbox: mbox,

		configStorage: configStorage,
		logStorage:    logStorage,
		logNotifier:   logNotifier,
	}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(mbox),
		proctree.NewProcess(handler),
	)

	return &Runner{
		mbox:    mbox,
		process: process,
	}
}

func (r *Runner) Start() {
	r.process.Start()
}

func (r *Runner) Stop() {
	r.process.Stop()
}

func (r *Runner) WaitReady(ctx context.Context) error {
	return r.process.WaitReady(ctx)
}

func (r *Runner) Join(ctx context.Context) error {
	return r.process.Join(ctx)
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
