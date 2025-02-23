package lognotify

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/models"
)

type LogNotifier struct {
	subMbox actor.MailboxSender[SubscribeMessage]
	logMbox actor.MailboxSender[LogMessage]
	process proctree.Process
}

func NewLogNotifier() *LogNotifier {
	subMbox := actor.NewMailbox[SubscribeMessage]()
	logMbox := actor.NewMailbox[LogMessage]()
	handler := &procHandler{
		subscribers: make(map[string]subscriberSet),
		subMbox:     subMbox,
		logMbox:     logMbox,
	}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(subMbox),
		proctree.NewActorProcess(logMbox),
		proctree.NewProcess(handler),
	)

	return &LogNotifier{
		subMbox: subMbox,
		logMbox: logMbox,
		process: process,
	}
}

func (n *LogNotifier) Start() {
	n.process.Start()
}

func (n *LogNotifier) Stop() {
	n.process.Stop()
}

func (n *LogNotifier) WaitReady(ctx context.Context) error {
	return n.process.WaitReady(ctx)
}

func (n *LogNotifier) Join(ctx context.Context) error {
	return n.process.Join(ctx)
}

func (n *LogNotifier) Subscribe(ctx context.Context, stream string) (actor.MailboxReceiver[LogMessage], error) {
	logM := actor.NewMailbox[LogMessage]()
	logM.Start()

	readyC := make(chan ReadyResponse, 1)
	msg := SubscribeMessage{
		Stream:  stream,
		SenderM: logM,
		ReadyC:  readyC,
		DoneC:   ctx.Done(),
	}

	err := n.subMbox.Send(ctx, msg)
	if err != nil {
		logM.Stop()
		return nil, err
	}

	resp := <-readyC
	if resp.Err != nil {
		logM.Stop()
		return nil, resp.Err
	}

	return logM, nil
}

func (n *LogNotifier) Notify(ctx context.Context, stream string, logKey string, logRecord models.LogRecord) error {
	msg := LogMessage{
		Stream:    stream,
		LogKey:    logKey,
		LogRecord: logRecord,
	}
	return n.logMbox.Send(ctx, msg)
}
