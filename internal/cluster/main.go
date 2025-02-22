package cluster

import (
	"log/slog"

	"context"

	"crypto/tls"
	"net"
	"net/http"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/cluster/rafthttp"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type Manager struct {
	state   *state
	process proctree.Process
}

func NewManager(
	dir string,
	nodeID string,
	joinAddr string,

	listener net.Listener,
	tlsConfig *tls.Config,
	timeout time.Duration,

	authStorage *auth.Storage,
	configStorage *config.Storage,
	logStorage *log.Storage,
) *Manager {
	var dial rafthttp.Dial
	if tlsConfig != nil {
		dial = rafthttp.NewDialTLS(tlsConfig)
	} else {
		dial = rafthttp.NewDialTCP()
	}

	connM := actor.NewMailbox[net.Conn]()

	state := &state{
		logger: slog.Default().With(
			slog.String("channel", "cluster"),
			slog.Group("raft",
				slog.String("dir", dir),
				slog.String("node-id", nodeID),
				slog.String("node-addr", listener.Addr().String()),
				slog.String("join-addr", joinAddr),
			),
		),

		listener: listener,
		dial:     dial,
		timeout:  timeout,
		connM:    connM,

		nodeID:   nodeID,
		joinAddr: joinAddr,

		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,
	}

	process := proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewActorProcess(connM),
		proctree.NewProcess(&storeHandler{state: state, dir: dir}),
		proctree.NewProcess(&snapshotsHandler{state: state, dir: dir}),
		proctree.NewProcess(&consensusHandler{state: state}),
		proctree.NewProcess(&membershipHandler{state: state}),
		proctree.NewProcess(&httpHandler{state: state}),
		proctree.NewProcess(&joinHandler{state: state}),
	)

	return &Manager{
		state:   state,
		process: process,
	}
}

func (m *Manager) Start() {
	m.process.Start()
}

func (m *Manager) Stop() {
	m.process.Stop()
}

func (m *Manager) WaitReady(ctx context.Context) error {
	return m.process.WaitReady(ctx)
}

func (m *Manager) Join(ctx context.Context) error {
	return m.process.Join(ctx)
}

func (m *Manager) HttpHandler() http.Handler {
	return m.state.httpHandler
}
