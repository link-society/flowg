package http

import (
	"context"
	"fmt"

	"log/slog"

	"crypto/tls"
	"net"
	"net/http"

	"go.uber.org/fx"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"
)

type ServerOptions struct {
	BindAddress string
	MountPath   string
	TlsConfig   *tls.Config
}

type Server struct {
	logger     *slog.Logger
	httpServer *http.Server
}

type handlers struct {
	fx.In

	ApiHandler http.Handler `name:"service-http-api"`
	WebHandler http.Handler `name:"service-http-web"`
}

func NewServer(opts ServerOptions) fx.Option {
	return fx.Module(
		"services.http",
		api.Module("service-http-api"),
		web.Module("service-http-web", opts.MountPath),
		fx.Provide(func(lc fx.Lifecycle, h handlers) *Server {
			rootHandler := http.NewServeMux()
			rootHandler.Handle(opts.MountPath+"/api/", http.StripPrefix(opts.MountPath, h.ApiHandler))
			rootHandler.Handle(opts.MountPath+"/web/", http.StripPrefix(opts.MountPath, h.WebHandler))

			rootHandler.HandleFunc(
				"GET "+opts.MountPath+"/{$}",
				func(w http.ResponseWriter, r *http.Request) {
					http.Redirect(w, r, opts.MountPath+"/web/", http.StatusPermanentRedirect)
				},
			)

			srv := &Server{
				logger: slog.Default().With(
					slog.String("channel", "http"),
					slog.Group("http",
						slog.String("bind", opts.BindAddress),
					),
				),
				httpServer: &http.Server{
					Addr:      opts.BindAddress,
					Handler:   newLoggingMiddleware(rootHandler),
					TLSConfig: opts.TlsConfig,
				},
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					srv.logger.InfoContext(ctx, "Starting HTTP server")

					l, err := net.Listen("tcp", opts.BindAddress)
					if err != nil {
						return fmt.Errorf("failed to start HTTP server: %w", err)
					}

					if opts.TlsConfig != nil {
						go srv.httpServer.ServeTLS(l, "", "")
					} else {
						go srv.httpServer.Serve(l)
					}

					return nil
				},
				OnStop: func(ctx context.Context) error {
					srv.logger.InfoContext(ctx, "Stopping HTTP server")
					return srv.httpServer.Shutdown(ctx)
				},
			})

			return srv
		}),
	)
}
