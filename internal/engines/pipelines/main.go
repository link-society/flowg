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
	proctree.Process

	mbox actor.MailboxSender[message]
}

var _ proctree.Process = (*Runner)(nil)

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
		Process: process,
		mbox:    mbox,
	}
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
