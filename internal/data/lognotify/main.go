package lognotify

import (
	"context"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/data/logstorage"
)

type LogNotifier struct {
	subMbox actor.MailboxSender[SubscribeMessage]
	logMbox actor.MailboxSender[LogMessage]
	rootA   actor.Actor
}

func NewLogNotifier() *LogNotifier {
	subMbox := actor.NewMailbox[SubscribeMessage]()
	logMbox := actor.NewMailbox[LogMessage]()
	workerA := actor.New(&worker{
		subscribers: make(map[string]subscriberSet),
		subMbox:     subMbox,
		logMbox:     logMbox,
	})

	return &LogNotifier{
		subMbox: subMbox,
		logMbox: logMbox,
		rootA:   actor.Combine(subMbox, logMbox, workerA).Build(),
	}
}

func (n *LogNotifier) Start() {
	n.rootA.Start()
}

func (n *LogNotifier) Stop() {
	n.rootA.Stop()
}

func (n *LogNotifier) Subscribe(ctx context.Context, stream string) actor.MailboxReceiver[LogMessage] {
	logM := actor.NewMailbox[LogMessage]()
	logM.Start()

	msg := SubscribeMessage{
		Stream:  stream,
		SenderM: logM,
		DoneC:   ctx.Done(),
	}

	n.subMbox.Send(ctx, msg)

	return logM
}

func (n *LogNotifier) Notify(ctx context.Context, stream string, logKey string, logEntry logstorage.LogEntry) {
	msg := LogMessage{
		Stream:   stream,
		LogKey:   logKey,
		LogEntry: logEntry,
	}
	n.logMbox.Send(ctx, msg)
}
