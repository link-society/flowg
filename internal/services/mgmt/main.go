package mgmt

import (
	"context"
	"log/slog"

	"crypto/tls"
	"net/http"

	"go.uber.org/fx"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/cluster"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterNodeID            string
	ClusterCookie            string
	ClusterStateDir          string
	ClusterFormationStrategy cluster.ClusterFormationStrategy
}

type Server struct {
	httpServer *http.Server
}

func NewServer(opts ServerOptions) fx.Option {
	logger := slog.Default().With(
		slog.String("channel", "mgmt"),
		slog.Group("mgmt",
			slog.String("bind", opts.BindAddress),
		),
	)

	return fx.Module(
		"services.mgmt",
		fx.Provide(func() (*cluster.Listener, error) {
			return cluster.NewListener(
				slog.Default().With(slog.String("channel", "cluster.listener")),
				opts.BindAddress,
				opts.TlsConfig,
			)
		}),
		cluster.NewManager(cluster.ManagerOptions{
			NodeID: opts.ClusterNodeID,
			Cookie: opts.ClusterCookie,

			ClusterFormationStrategy: opts.ClusterFormationStrategy,
			ClusterStateDir:          opts.ClusterStateDir,
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			listener *cluster.Listener,
			manager cluster.Manager,
		) *Server {
			srv := &Server{}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Start Management HTTP server")

					metricsHandler := promhttp.Handler()
					clusterHandler := manager.HttpHandler()

					rootHandler := http.NewServeMux()
					registerProfiler(rootHandler)

					rootHandler.HandleFunc(
						"/health",
						func(w http.ResponseWriter, r *http.Request) {
							w.WriteHeader(http.StatusOK)
							w.Write([]byte("OK\r\n"))
						},
					)
					rootHandler.Handle("/metrics", metricsHandler)
					rootHandler.Handle("/cluster/", clusterHandler)

					srv.httpServer = &http.Server{
						Addr:      opts.BindAddress,
						Handler:   rootHandler,
						TLSConfig: opts.TlsConfig,
					}

					if opts.TlsConfig != nil {
						go srv.httpServer.ServeTLS(listener.Socket(), "", "")
					} else {
						go srv.httpServer.Serve(listener.Socket())
					}

					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Stopping Management HTTP server")
					return srv.httpServer.Shutdown(ctx)
				},
			})

			return srv
		}),
	)
}
