package server

import (
	"log/slog"
	"time"

	"crypto/tls"

	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

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

	ClusterNodeID            string
	ClusterCookie            string
	ClusterStateDir          string
	ClusterFormationStrategy cluster.ClusterFormationStrategy

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

	AuthResetUser     string
	AuthResetPassword string
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

			ClusterNodeID:            opts.ClusterNodeID,
			ClusterCookie:            opts.ClusterCookie,
			ClusterStateDir:          opts.ClusterStateDir,
			ClusterFormationStrategy: opts.ClusterFormationStrategy,

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
		resetUser:       opts.AuthResetUser,
		resetPassword:   opts.AuthResetPassword,
	}

	return proctree.NewProcessGroup(
		proctree.ProcessGroupOptions{
			// Longer init timeout because discovering other nodes
			// could take longer than the default 5 seconds
			InitTimeout: 1 * time.Minute,
			JoinTimeout: 5 * time.Second,
		},
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
			proctree.ProcessGroupOptions{
				// Longer init timeout because discovering other nodes
				// could take longer than the default 5 seconds
				InitTimeout: 1 * time.Minute,
				JoinTimeout: 5 * time.Second,
			},
			httpServer,
			mgmtServer,
			syslogServer,
		),
		proctree.NewProcess(bootstrapProc),
	)
}
