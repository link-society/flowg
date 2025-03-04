package syslog

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/engines/pipelines"
)

type ServerOptions struct {
	TcpMode      bool
	BindAddress  string
	TlsConfig    *tls.Config
	AllowOrigins []string

	ConfigStorage  *config.Storage
	PipelineRunner *pipelines.Runner
}

func NewServer(opts *ServerOptions) proctree.Process {
	proto := "udp"
	if opts.TcpMode {
		proto = "tcp"
	}

	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "syslog"),
			slog.Group("syslog",
				slog.String("proto", proto),
				slog.String("bind", opts.BindAddress),
				slog.Bool("tls", opts.TlsConfig != nil),
			),
		),

		opts: opts,
	})
}
