package server

import (
	"context"
	"log/slog"

	"crypto/tls"

	"go.uber.org/fx"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
)

// Options configures the FlowG server: the bind addresses and TLS settings of
// each network service, the on-disk locations of the three storage backends,
// and the initial/reset credentials applied during bootstrap.
type Options struct {
	HttpBindAddress string
	HttpMountPath   string
	HttpTlsConfig   *tls.Config

	MgmtBindAddress string
	MgmtTlsConfig   *tls.Config

	SyslogTcpMode               bool
	SyslogBindAddress           string
	SyslogTlsConfig             *tls.Config
	SyslogInitialAllowedOrigins []string

	Storage StorageOptions

	AuthInitialUser     string
	AuthInitialPassword string

	AuthResetUser     string
	AuthResetPassword string
}

// StorageOptions is an interface that abstracts the storage backend
// configuration.
type StorageOptions interface {
	// AuthModule returns an fx module that provides the auth storage backend.
	AuthModule() fx.Option
	// ConfigModule returns an fx module that provides the config storage backend.
	ConfigModule() fx.Option
	// LogModule returns an fx module that provides the log storage backend.
	LogModule() fx.Option
}

// NewServer assembles the complete FlowG server as a single fx module. It wires
// the three storage backends, the engine layer (log notifier and pipeline
// runner) and the network services (HTTP, management, syslog) together, and
// registers a bootstrap handler that seeds the default configuration on start.
func NewServer(opts Options) fx.Option {
	return fx.Module(
		"app.server",
		// Storage Layer
		opts.Storage.AuthModule(),
		opts.Storage.ConfigModule(),
		opts.Storage.LogModule(),
		// Engine Layer
		lognotify.NewLogNotifier(),
		pipelines.NewRunner(),
		// Service Layer
		http.NewServer(http.ServerOptions{
			BindAddress: opts.HttpBindAddress,
			MountPath:   opts.HttpMountPath,
			TlsConfig:   opts.HttpTlsConfig,
		}),
		mgmt.NewServer(mgmt.ServerOptions{
			BindAddress: opts.MgmtBindAddress,
			TlsConfig:   opts.MgmtTlsConfig,
		}),
		syslog.NewServer(syslog.ServerOptions{
			TcpMode:     opts.SyslogTcpMode,
			BindAddress: opts.SyslogBindAddress,
			TlsConfig:   opts.SyslogTlsConfig,
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			authStorage storage.AuthStorage,
			configStorage storage.ConfigStorage,
			logStorage storage.LogStorage,
		) *bootstrapHandler {
			h := &bootstrapHandler{
				logger:                      slog.Default().With(slog.String("channel", "bootstrap")),
				authStorage:                 authStorage,
				configStorage:               configStorage,
				initialSyslogAllowedOrigins: opts.SyslogInitialAllowedOrigins,
				initialUser:                 opts.AuthInitialUser,
				initialPassword:             opts.AuthInitialPassword,
				resetUser:                   opts.AuthResetUser,
				resetPassword:               opts.AuthResetPassword,
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return h.Run(ctx)
				},
			})

			return h
		}),
		fx.Invoke(func(_ struct {
			fx.In

			// Storage layer
			AuthStorage   storage.AuthStorage
			ConfigStorage storage.ConfigStorage
			LogStorage    storage.LogStorage
			// Engine layer
			LogNotifier     lognotify.LogNotifier
			PipelinesRunner pipelines.Runner
			// Service layer
			HttpServer       *http.Server
			ManagementServer *mgmt.Server
			SyslogServer     *syslog.Server
			// Bootstrap handler
			BootstrapHandler *bootstrapHandler
		}) {
			// No-op, just to force the creation of all components
		}),
	)
}
