package syslog

import (
	"context"
	"fmt"
	"log/slog"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	storage "link-society.com/flowg/internal/storage/interfaces"

	"link-society.com/flowg/internal/engines/pipelines"
)

// ServerOptions configures the syslog server: UDP (default) or TCP, where to
// bind, and an optional TLS configuration (TCP only).
type ServerOptions struct {
	TcpMode     bool
	BindAddress string
	TlsConfig   *tls.Config
}

// Server is the running syslog service: an actor that drains received messages
// into the pipeline engine.
type Server struct {
	actor.Actor
}

// NewServer returns an fx module that listens for syslog messages (auto-detecting
// the format) and feeds each one, through the worker actor, into every pipeline's
// syslog entrypoint. Listener and actor are bound to the application lifecycle.
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
			configStorage storage.ConfigStorage,
			pipelineRunner pipelines.Runner,
		) *Server {
			srv := &Server{
				Actor: actor.New(&worker{
					logger:         logger,
					channel:        channel,
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
		fx.Invoke(func(*gosyslog.Server) {
			// No-op, just to force the creation of all components
		}),
	)
}
