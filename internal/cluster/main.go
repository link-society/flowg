package cluster

import (
	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/memberlist"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type ManagerOptions struct {
	NodeID string
	Cookie string

	ClusterJoinNode *ClusterJoinNode

	AutomaticClusterFormation bool

	LocalEndpointResolver func() (*url.URL, error)

	AuthStorage   *auth.Storage
	ConfigStorage *config.Storage
	LogStorage    *log.Storage
}

type Manager struct {
	proctree.Process

	handler *procHandler
}

var _ proctree.Process = (*Manager)(nil)

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
