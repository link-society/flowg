package http

import (
	"context"
	"fmt"

	"log/slog"

	"crypto/tls"
	"net"
	gohttp "net/http"

	"go.uber.org/fx"

	"link-society.com/flowg/api"
	_ "link-society.com/flowg/api/operations"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/web"

	"link-society.com/flowg/internal/app/logging"
)

type ServerOptions struct {
	BindAddress string
	MountPath   string
	TlsConfig   *tls.Config
}

type Server struct {
	logger     *slog.Logger
	httpServer *gohttp.Server
}

type handlers struct {
	fx.In

	ApiHandler gohttp.Handler `name:"service-http-api"`
	WebHandler gohttp.Handler `name:"service-http-web"`
}

func NewServer(opts ServerOptions) fx.Option {
	return fx.Module(
		"services.http",
		routing.Module("api.operations"),
		fx.Provide(fx.Annotate(
			api.NewHandler,
			fx.ResultTags(`name:"service-http-api"`),
		)),
		fx.Provide(fx.Annotate(
			func() gohttp.Handler {
				return web.NewHandler(opts.MountPath)
			},
			fx.ResultTags(`name:"service-http-web"`),
		)),
		fx.Provide(func(lc fx.Lifecycle, h handlers) *Server {
			rootHandler := gohttp.NewServeMux()
			rootHandler.Handle(opts.MountPath+"/api/", gohttp.StripPrefix(opts.MountPath, h.ApiHandler))
			rootHandler.Handle(opts.MountPath+"/web/", gohttp.StripPrefix(opts.MountPath, h.WebHandler))

			rootHandler.HandleFunc(
				"GET "+opts.MountPath+"/{$}",
				func(w gohttp.ResponseWriter, r *gohttp.Request) {
					gohttp.Redirect(w, r, opts.MountPath+"/web/", gohttp.StatusPermanentRedirect)
				},
			)

			srv := &Server{
				logger: slog.Default().With(
					slog.String("channel", "http"),
					slog.Group("http",
						slog.String("bind", opts.BindAddress),
					),
				),
				httpServer: &gohttp.Server{
					Addr:      opts.BindAddress,
					Handler:   logging.NewMiddleware(rootHandler),
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
