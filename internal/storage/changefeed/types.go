package changefeed

import (
	"github.com/vladopajic/go-actor/actor"
)

type Operation int

const (
	OpWrite Operation = iota
	OpDelete
)

const (
	NamespaceConfig = "config"
	NamespaceAuth   = "auth"
	NamespaceLog    = "log"
)

type ChangeEvent struct {
	Namespace string
	Kind      string
	Name      string
	Op        Operation
	Resync    bool
}

type subscribeMessage struct {
	senderM actor.Mailbox[ChangeEvent]
	readyC  chan<- error
	doneC   <-chan struct{}
}
