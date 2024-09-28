package lognotify

import (
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"
)

type LogMessage struct {
	Stream    string
	LogKey    string
	LogRecord models.LogRecord
}

type SubscribeMessage struct {
	Stream  string
	SenderM actor.Mailbox[LogMessage]
	ReadyC  chan<- ReadyResponse
	DoneC   <-chan struct{}
}

type ReadyResponse struct {
	Err error
}
