package server

import (
	"log/slog"
	"net/url"

	"crypto/tls"

	"link-society.com/flowg/internal/cluster"
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
	ClusterJoinNodeID   string
	ClusterJoinEndpoint *url.URL
	ClusterCookie       string
	ClusterStateDir     string

	SyslogTcpMode      bool
	SyslogBindAddress  string
	SyslogTlsConfig    *tls.Config
	SyslogAllowOrigins []string

	AuthStorageDir   string
	ConfigStorageDir string
	LogStorageDir    string

	ServiceName string
	ConsulUrl   string

	AuthInitialUser     string
	AuthInitialPassword string
}

func NewServer(opts Options) proctree.Process {
	isAutomaticClusterFormation := isAutomaticClusterFormation(opts.ConsulUrl)

	// ClusterJoinNode shared between ConsulService and ManagementServer
	ClusterJoinNode := cluster.NewClusterJoinNode(isAutomaticClusterFormation, opts.ClusterJoinNodeID, opts.ClusterJoinEndpoint)

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
		consulService = consul.NewConsulService(&consul.ConsulServiceOptions{
			NodeId:          opts.ClusterNodeID,
			ServiceName:     opts.ServiceName,
			ConsulUrl:       opts.ConsulUrl,
			ClusterJoinNode: ClusterJoinNode,
			MgmtBindAddress: opts.MgmtBindAddress,
			MgmtTlsEnabled:  opts.MgmtTlsConfig != nil,
		})
		mgmtServer = mgmt.NewServer(&mgmt.ServerOptions{
			BindAddress: opts.MgmtBindAddress,
			TlsConfig:   opts.MgmtTlsConfig,

			ClusterNodeID:   opts.ClusterNodeID,
			ClusterCookie:   opts.ClusterCookie,
			ClusterJoinNode: ClusterJoinNode,
			ClusterStateDir: opts.ClusterStateDir,

			AutomaticClusterFormation: isAutomaticClusterFormation,

			AuthStorage:   authStorage,
			ConfigStorage: configStorage,
			LogStorage:    logStorage,
		})
		syslogServer = syslog.NewServer(&syslog.ServerOptions{
			TcpMode:      opts.SyslogTcpMode,
			BindAddress:  opts.SyslogBindAddress,
			TlsConfig:    opts.SyslogTlsConfig,
			AllowOrigins: opts.SyslogAllowOrigins,

			ConfigStorage:  configStorage,
			PipelineRunner: pipelineRunner,
		})
	)

	// Bootstrap Process
	bootstrapProc := &bootstrapProcHandler{
		logger:          slog.Default().With(slog.String("channel", "bootstrap")),
		authStorage:     authStorage,
		configStorage:   configStorage,
		initialUser:     opts.AuthInitialUser,
		initialPassword: opts.AuthInitialPassword,
	}

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
			consulService,
			mgmtServer,
			syslogServer,
		),
		proctree.NewProcess(bootstrapProc),
	)
}

func isAutomaticClusterFormation(consulUrl string) bool {
	return consulUrl != ""
}
