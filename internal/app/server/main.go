package server

import (
	"context"
	"log/slog"

	"crypto/tls"

	"go.uber.org/fx"

	"link-society.com/flowg/internal/cluster"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
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

func NewServer(opts Options) fx.Option {
	return fx.Module(
		"app.server",
		// Storage Layer
		auth.NewStorage(func() auth.Options {
			authOpts := auth.DefaultOptions()
			authOpts.Directory = opts.AuthStorageDir
			return authOpts
		}()),
		config.NewStorage(func() config.Options {
			configOpts := config.DefaultOptions()
			configOpts.Directory = opts.ConfigStorageDir
			return configOpts
		}()),
		log.NewStorage(func() log.Options {
			logOpts := log.DefaultOptions()
			logOpts.Directory = opts.LogStorageDir
			return logOpts
		}()),
		// Engine Layer
		lognotify.NewLogNotifier(),
		pipelines.NewRunner(),
		// Service Layer
		http.NewServer(http.ServerOptions{
			BindAddress: opts.HttpBindAddress,
			TlsConfig:   opts.HttpTlsConfig,
		}),
		mgmt.NewServer(mgmt.ServerOptions{
			BindAddress: opts.MgmtBindAddress,
			TlsConfig:   opts.MgmtTlsConfig,

			ClusterNodeID:            opts.ClusterNodeID,
			ClusterCookie:            opts.ClusterCookie,
			ClusterStateDir:          opts.ClusterStateDir,
			ClusterFormationStrategy: opts.ClusterFormationStrategy,
		}),
		syslog.NewServer(syslog.ServerOptions{
			TcpMode:      opts.SyslogTcpMode,
			BindAddress:  opts.SyslogBindAddress,
			TlsConfig:    opts.SyslogTlsConfig,
			AllowOrigins: opts.SyslogAllowOrigins,
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			authStorage auth.Storage,
			configStorage config.Storage,
			logStorage log.Storage,
		) *bootstrapHandler {
			h := &bootstrapHandler{
				logger:          slog.Default().With(slog.String("channel", "bootstrap")),
				authStorage:     authStorage,
				configStorage:   configStorage,
				initialUser:     opts.AuthInitialUser,
				initialPassword: opts.AuthInitialPassword,
				resetUser:       opts.AuthResetUser,
				resetPassword:   opts.AuthResetPassword,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return h.Run(ctx)
				},
			})

			return h
		}),
		fx.Invoke(func(
			// Storage layer
			_ auth.Storage,
			_ config.Storage,
			_ log.Storage,
			// Engine layer
			_ lognotify.LogNotifier,
			_ pipelines.Runner,
			// Service layer
			_ *http.Server,
			_ *mgmt.Server,
			_ *syslog.Server,
			// Bootstrap handler
			_ *bootstrapHandler,
		) {
		}),
	)
}
