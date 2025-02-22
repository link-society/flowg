package mgmt

import (
	"log/slog"

	"time"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterStateDir string
	ClusterNodeID   string
	ClusterJoinAddr string
	ClusterTimeout  time.Duration

	AuthStorage   *auth.Storage
	ConfigStorage *config.Storage
	LogStorage    *log.Storage
}

func NewServer(opts ServerOptions) proctree.Process {
	state := &state{
		logger: slog.Default().With(
			slog.String("channel", "mgmt"),
			slog.Group("mgmt",
				slog.String("bind", opts.BindAddress),
				slog.Bool("tls", opts.TlsConfig != nil),
			),
			slog.Group("cluster",
				slog.String("dir", opts.ClusterStateDir),
				slog.String("node-id", opts.ClusterNodeID),
				slog.String("join-addr", opts.ClusterJoinAddr),
			),
		),

		bindAddress: opts.BindAddress,
		tlsConfig:   opts.TlsConfig,

		clusterStateDir: opts.ClusterStateDir,
		clusterNodeID:   opts.ClusterNodeID,
		clusterJoinAddr: opts.ClusterJoinAddr,
		clusterTimeout:  opts.ClusterTimeout,

		authStorage:   opts.AuthStorage,
		configStorage: opts.ConfigStorage,
		logStorage:    opts.LogStorage,
	}

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewProcess(&listenerHandler{state: state}),
		proctree.NewProcess(&clusterManagerHandler{state: state}),
		proctree.NewProcess(&httpHandler{state: state}),
	)
}
