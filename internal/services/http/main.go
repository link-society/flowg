package http

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	AuthStorage   auth.Storage
	ConfigStorage config.Storage
	LogStorage    log.Storage

	LogNotifier    *lognotify.LogNotifier
	PipelineRunner *pipelines.Runner
}

func NewServer(opts *ServerOptions) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "http"),
			slog.Group("http",
				slog.String("bind", opts.BindAddress),
			),
		),

		opts: opts,
	})
}
