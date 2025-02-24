package server

import (
	"time"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
)

type serviceLayerOpts struct {
	httpBindAddress string
	httpTlsConfig   *tls.Config

	mgmtBindAddress string
	mgmtTlsConfig   *tls.Config

	clusterStateDir string
	clusterNodeID   string
	clusterJoinAddr string
	clusterTimeout  time.Duration

	syslogTCP          bool
	syslogBindAddress  string
	syslogTlsConfig    *tls.Config
	syslogAllowOrigins []string

	storageLayer *storageLayer
	engineLayer  *engineLayer
}

func newServiceLayer(opts serviceLayerOpts) proctree.Process {
	httpServer := http.NewServer(
		opts.httpBindAddress,
		opts.httpTlsConfig,
		opts.storageLayer.authStorage,
		opts.storageLayer.configStorage,
		opts.storageLayer.logStorage,
		opts.engineLayer.logNotifier,
		opts.engineLayer.pipelineRunner,
	)

	mgmtServer := mgmt.NewServer(mgmt.ServerOptions{
		BindAddress: opts.mgmtBindAddress,
		TlsConfig:   opts.mgmtTlsConfig,

		ClusterStateDir: opts.clusterStateDir,
		ClusterNodeID:   opts.clusterNodeID,
		ClusterJoinAddr: opts.clusterJoinAddr,
		ClusterTimeout:  opts.clusterTimeout,

		AuthStorage:   opts.storageLayer.authStorage,
		ConfigStorage: opts.storageLayer.configStorage,
		LogStorage:    opts.storageLayer.logStorage,
	})

	syslogServer := syslog.NewServer(
		opts.syslogTCP,
		opts.syslogBindAddress,
		opts.syslogTlsConfig,
		opts.syslogAllowOrigins,

		opts.storageLayer.configStorage,
		opts.engineLayer.pipelineRunner,
	)

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		httpServer,
		mgmtServer,
		syslogServer,
	)
}
