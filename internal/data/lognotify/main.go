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
		subscribers: make(map[string]map[chan<- LogMessage]struct{}),
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

func (n *LogNotifier) Subscribe(stream string, doneC <-chan struct{}) <-chan LogMessage {
	logC := make(chan LogMessage)
	msg := SubscribeMessage{
		Stream:  stream,
		SenderC: logC,
		DoneC:   doneC,
	}

	n.subMbox.Send(context.Background(), msg)

	return logC
}

func (n *LogNotifier) Notify(stream string, logKey string, logEntry logstorage.LogEntry) {
	msg := LogMessage{
		Stream:   stream,
		LogKey:   logKey,
		LogEntry: logEntry,
	}
	n.logMbox.Send(context.Background(), msg)
}
