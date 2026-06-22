package mgmt

import (
	"context"
	"fmt"
	"log/slog"

	"crypto/tls"
	"net"
	"net/http"

	"go.uber.org/fx"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config
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
		fx.Provide(func(lc fx.Lifecycle) *Server {
			srv := &Server{}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Start Management HTTP server")

					metricsHandler := promhttp.Handler()

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

					l, err := net.Listen("tcp", opts.BindAddress)
					if err != nil {
						return fmt.Errorf("failed to start Management HTTP server: %w", err)
					}

					srv.httpServer = &http.Server{
						Addr:      opts.BindAddress,
						Handler:   rootHandler,
						TLSConfig: opts.TlsConfig,
					}

					if opts.TlsConfig != nil {
						go srv.httpServer.ServeTLS(l, "", "")
					} else {
						go srv.httpServer.Serve(l)
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
