package cluster

import (
	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/memberlist"
)

type Manager struct {
	proctree.Process

	handler *procHandler
}

func NewManager(
	nodeID string,
	localEndpoint *url.URL,
	joinNodeID string,
	joinNodeEndpoint *url.URL,
) *Manager {
	connM := actor.NewMailbox[net.Conn]()
	packetM := actor.NewMailbox[*memberlist.Packet]()
	handler := &procHandler{
		nodeID:           nodeID,
		localEndpoint:    localEndpoint,
		joinNodeID:       joinNodeID,
		joinNodeEndpoint: joinNodeEndpoint,

		connM:   connM,
		packetM: packetM,
	}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(connM),
		proctree.NewActorProcess(packetM),
		proctree.NewProcess(handler),
	)

	return &Manager{
		Process: process,
		handler: handler,
	}
}

func (m *Manager) HttpHandler() http.Handler {
	return m.handler.httpHandler
}
