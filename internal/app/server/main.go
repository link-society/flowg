package server

import (
	"log/slog"

	"go.uber.org/fx"

	"link-society.com/flowg/internal/app/config"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/confignotify"
	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
)

// NewServer assembles the complete FlowG server as a single fx module. It
// supplies the configuration and the root logger to the graph, wires the three storage backends,
// the engine layer (log notifier and pipeline runner) and the network services
// (HTTP, management, syslog) together, and registers a bootstrap handler that
// seeds the default configuration on start.
func NewServer(cfg *config.RootConfig, logger *slog.Logger) fx.Option {
	return fx.Module(
		"app.server",
		// Configuration
		configModule(cfg),
		fx.Supply(logger),
		// Engine Layer (confignotify first: storage modules depend on it)
		confignotify.NewNotifier(),
		// Storage Layer
		storageModules(cfg),
		// Engine Layer
		lognotify.NewLogNotifier(),
		pipelines.NewRunner(),
		// Service Layer
		http.Module(),
		mgmt.Module(),
		syslog.Module(),
		// Bootstrap
		fx.Provide(newBootstrapHandler),
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
