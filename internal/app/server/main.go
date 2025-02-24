package server

import (
	"log/slog"

	"time"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"
)

type Options struct {
	HttpBindAddress string
	HttpTlsConfig   *tls.Config

	MgmtBindAddress string
	MgmtTlsConfig   *tls.Config

	ClusterStateDir string
	ClusterNodeID   string
	ClusterJoinAddr string

	SyslogTCP          bool
	SyslogBindAddress  string
	SyslogTlsConfig    *tls.Config
	SyslogAllowOrigins []string

	AuthStorageDir   string
	ConfigStorageDir string
	LogStorageDir    string
}

func NewServer(opts Options) proctree.Process {
	storageLayer := newStorageLayer(
		opts.AuthStorageDir,
		opts.ConfigStorageDir,
		opts.LogStorageDir,
	)
	engineLayer := newEngineLayer(
		storageLayer,
	)
	serviceLayer := newServiceLayer(serviceLayerOpts{
		httpBindAddress: opts.HttpBindAddress,
		httpTlsConfig:   opts.HttpTlsConfig,

		mgmtBindAddress: opts.MgmtBindAddress,
		mgmtTlsConfig:   opts.MgmtTlsConfig,

		clusterStateDir: opts.ClusterStateDir,
		clusterNodeID:   opts.ClusterNodeID,
		clusterJoinAddr: opts.ClusterJoinAddr,
		clusterTimeout:  5 * time.Second,

		syslogTCP:          opts.SyslogTCP,
		syslogBindAddress:  opts.SyslogBindAddress,
		syslogTlsConfig:    opts.SyslogTlsConfig,
		syslogAllowOrigins: opts.SyslogAllowOrigins,

		storageLayer: storageLayer,
		engineLayer:  engineLayer,
	})

	bootstrap := proctree.NewProcess(&bootstrapProcHandler{
		logger:        slog.Default().With("channel", "server"),
		storageLayer:  storageLayer,
		isInitialNode: opts.ClusterJoinAddr == "",
	})

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		storageLayer,
		engineLayer,
		serviceLayer,
		bootstrap,
	)
}
