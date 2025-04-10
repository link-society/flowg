package server

import (
	"log/slog"
	"net/url"

	"crypto/tls"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/services/consul"
	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"

	"link-society.com/flowg/internal/utils/proctree"
)

type Options struct {
	HttpBindAddress string
	HttpTlsConfig   *tls.Config

	MgmtBindAddress string
	MgmtTlsConfig   *tls.Config

	ClusterNodeID       string
	ClusterNodeAddress  *url.URL
	ClusterJoinNodeID   string
	ClusterJoinEndpoint *url.URL
	ClusterCookie       string

	SyslogTcpMode      bool
	SyslogBindAddress  string
	SyslogTlsConfig    *tls.Config
	SyslogAllowOrigins []string

	AuthStorageDir   string
	ConfigStorageDir string
	LogStorageDir    string

	ServiceName string
	ConsulUrl   string
}

func NewServer(opts Options) proctree.Process {
	// Storage Layer
	var (
		authStorage   = auth.NewStorage(auth.OptDirectory(opts.AuthStorageDir))
		configStorage = config.NewStorage(config.OptDirectory(opts.ConfigStorageDir))
		logStorage    = log.NewStorage(log.OptDirectory(opts.LogStorageDir))
	)

	// Engine Layer
	var (
		logNotifier    = lognotify.NewLogNotifier()
		pipelineRunner = pipelines.NewRunner(configStorage, logStorage, logNotifier)
	)

	// Service Layer
	var (
		httpServer = http.NewServer(&http.ServerOptions{
			BindAddress: opts.HttpBindAddress,
			TlsConfig:   opts.HttpTlsConfig,

			AuthStorage:   authStorage,
			ConfigStorage: configStorage,
			LogStorage:    logStorage,

			LogNotifier:    logNotifier,
			PipelineRunner: pipelineRunner,
		})
		mgmtServer = mgmt.NewServer(&mgmt.ServerOptions{
			BindAddress: opts.MgmtBindAddress,
			TlsConfig:   opts.MgmtTlsConfig,

			ClusterNodeID:       opts.ClusterNodeID,
			ClusterJoinNodeID:   opts.ClusterJoinNodeID,
			ClusterJoinEndpoint: opts.ClusterJoinEndpoint,
			ClusterCookie:       opts.ClusterCookie,
		})
		syslogServer = syslog.NewServer(&syslog.ServerOptions{
			TcpMode:      opts.SyslogTcpMode,
			BindAddress:  opts.SyslogBindAddress,
			TlsConfig:    opts.SyslogTlsConfig,
			AllowOrigins: opts.SyslogAllowOrigins,

			ConfigStorage:  configStorage,
			PipelineRunner: pipelineRunner,
		})

		consulService = consul.NewConsulService(&consul.ConsulServiceOptions{
			NodeId:      opts.ClusterNodeID,
			NodeHost:    opts.ClusterNodeAddress.Host,
			NodePort:    opts.ClusterNodeAddress.Port(),
			ServiceName: opts.ServiceName,
			ConsulUrl:   opts.ConsulUrl,
		})
	)

	bootstrap := proctree.NewProcess(&bootstrapProcHandler{
		logger: slog.Default().With("channel", "server"),

		authStorage:   authStorage,
		configStorage: configStorage,
	})

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewProcessGroup(
			proctree.DefaultProcessGroupOptions(),
			authStorage,
			configStorage,
			logStorage,
		),
		proctree.NewProcessGroup(
			proctree.DefaultProcessGroupOptions(),
			logNotifier,
			pipelineRunner,
		),
		proctree.NewProcessGroup(
			proctree.DefaultProcessGroupOptions(),
			httpServer,
			mgmtServer,
			syslogServer,
			consulService,
		),
		bootstrap,
	)
}
