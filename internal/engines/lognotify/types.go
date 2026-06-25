package lognotify

import (
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/models"
)

// LogMessage is a single record delivered to subscribers, carrying both the
// record and the storage key it was persisted under.
type LogMessage struct {
	Stream    string
	LogKey    string
	LogRecord models.LogRecord
}

// SubscribeMessage requests a new subscription on the actor. SenderM is the
// mailbox the actor should push matching records to, ReadyC reports when the
// registration is complete, and DoneC signals that the subscription should be
// dropped.
type SubscribeMessage struct {
	Stream  string
	SenderM actor.Mailbox[LogMessage]
	ReadyC  chan<- ReadyResponse
	DoneC   <-chan struct{}
}

// ReadyResponse acknowledges a SubscribeMessage, reporting any registration
// error.
type ReadyResponse struct {
	Err error
}
