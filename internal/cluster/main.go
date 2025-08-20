package cluster

import (
	"log/slog"

	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/kvstore"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/memberlist"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type Manager interface {
	proctree.Process

	HttpHandler() http.Handler
}

type ManagerOptions struct {
	NodeID string
	Cookie string

	ClusterFormationStrategy ClusterFormationStrategy

	LocalEndpointResolver func() (*url.URL, error)

	AuthStorage         auth.Storage
	ConfigStorage       config.Storage
	LogStorage          log.Storage
	ClusterStateStorage kvstore.Storage
}

type managerImpl struct {
	proctree.Process

	handler *procHandler
}

var _ Manager = (*managerImpl)(nil)

func NewManager(opts *ManagerOptions) Manager {
	connM := actor.NewMailbox[net.Conn]()
	packetM := actor.NewMailbox[*memberlist.Packet]()
	joinM := actor.NewMailbox[*ClusterJoinNode]()

	handler := &procHandler{
		opts: opts,

		connM:   connM,
		packetM: packetM,
		joinM:   joinM,
	}

	formationController := actor.New(&clusterFormationController{
		logger:   slog.Default().With(slog.String("channel", "cluster.formation")),
		joinM:    joinM,
		resolver: opts.LocalEndpointResolver,
		strategy: opts.ClusterFormationStrategy,
	})

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(connM),
		proctree.NewActorProcess(packetM),
		proctree.NewActorProcess(joinM),
		proctree.NewActorProcess(formationController),
		proctree.NewProcess(handler),
	)

	return &managerImpl{
		Process: process,
		handler: handler,
	}
}

func (m *managerImpl) HttpHandler() http.Handler {
	return m.handler.httpHandler
}
