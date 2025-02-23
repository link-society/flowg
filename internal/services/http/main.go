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

func NewServer(
	bindAddress string,
	tlsConfig *tls.Config,

	authStorage *auth.Storage,
	configStorage *config.Storage,
	logStorage *log.Storage,

	logNotifier *lognotify.LogNotifier,
	pipelineRunner *pipelines.Runner,
) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "http"),
			slog.Group("http",
				slog.String("bind", bindAddress),
			),
		),

		bindAddress: bindAddress,
		tlsConfig:   tlsConfig,

		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,

		logNotifier:    logNotifier,
		pipelineRunner: pipelineRunner,
	})
}
