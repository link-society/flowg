package syslog

import (
	"context"
	"fmt"
	"log/slog"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/engines/pipelines"
)

type ServerOptions struct {
	TcpMode      bool
	BindAddress  string
	TlsConfig    *tls.Config
	AllowOrigins []string
}

type Server struct {
	actor.Actor
}

func NewServer(opts ServerOptions) fx.Option {
	proto := "udp"
	if opts.TcpMode {
		proto = "tcp"
	}

	logger := slog.Default().With(
		slog.String("channel", "syslog"),
		slog.Group("syslog",
			slog.String("proto", proto),
			slog.String("bind", opts.BindAddress),
			slog.Bool("tls", opts.TlsConfig != nil),
		),
	)

	return fx.Module(
		"services.syslog",
		fx.Provide(func() gosyslog.LogPartsChannel {
			return make(gosyslog.LogPartsChannel)
		}),
		fx.Provide(func(channel gosyslog.LogPartsChannel) *gosyslog.ChannelHandler {
			return gosyslog.NewChannelHandler(channel)
		}),
		fx.Provide(func(lc fx.Lifecycle, handler *gosyslog.ChannelHandler) *gosyslog.Server {
			server := gosyslog.NewServer()
			server.SetFormat(gosyslog.Automatic)
			server.SetHandler(handler)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Starting Syslog server")

					switch {
					case opts.TcpMode && opts.TlsConfig != nil:
						if err := server.ListenTCPTLS(opts.BindAddress, opts.TlsConfig); err != nil {
							return fmt.Errorf("failed to listen on TCP+TLS: %w", err)
						}

					case opts.TcpMode && opts.TlsConfig == nil:
						if err := server.ListenTCP(opts.BindAddress); err != nil {
							return fmt.Errorf("failed to listen on TCP: %w", err)
						}

					case !opts.TcpMode:
						if err := server.ListenUDP(opts.BindAddress); err != nil {
							return fmt.Errorf("failed to listen on UDP: %w", err)
						}
					}

					if err := server.Boot(); err != nil {
						return fmt.Errorf("failed to boot server: %w", err)
					}

					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.InfoContext(ctx, "Stopping Syslog server")
					return server.Kill()
				},
			})

			return server
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			channel gosyslog.LogPartsChannel,
			configStorage config.Storage,
			pipelineRunner pipelines.Runner,
		) *Server {
			srv := &Server{
				Actor: actor.New(&worker{
					logger:         logger,
					channel:        channel,
					allowOrigins:   opts.AllowOrigins,
					configStorage:  configStorage,
					pipelineRunner: pipelineRunner,
				}),
			}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					srv.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					srv.Stop()
					return nil
				},
			})

			return srv
		}),
		fx.Invoke(func(*gosyslog.Server) {}),
	)
}
