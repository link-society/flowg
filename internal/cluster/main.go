package cluster

import (
	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/memberlist"
)

type ManagerOptions struct {
	NodeID           string
	JoinNodeID       string
	JoinNodeEndpoint *url.URL
	Cookie           string

	LocalEndpointResolver func() *url.URL
}

type Manager struct {
	proctree.Process

	handler *procHandler
}

func NewManager(opts *ManagerOptions) *Manager {
	connM := actor.NewMailbox[net.Conn]()
	packetM := actor.NewMailbox[*memberlist.Packet]()
	handler := &procHandler{
		opts: opts,

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
