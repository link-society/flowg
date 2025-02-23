package syslog

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/engines/pipelines"
)

func NewServer(
	isTCP bool,
	bindAddress string,
	tlsConfig *tls.Config,
	allowOrigins []string,

	configStorage *config.Storage,
	pipelineRunner *pipelines.Runner,
) proctree.Process {
	proto := "udp"
	if isTCP {
		proto = "tcp"
	}

	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "syslog"),
			slog.Group("syslog",
				slog.String("proto", proto),
				slog.String("bind", bindAddress),
				slog.Bool("tls", tlsConfig != nil),
			),
		),

		isTCP:        isTCP,
		bindAddress:  bindAddress,
		tlsConfig:    tlsConfig,
		allowOrigins: allowOrigins,

		configStorage:  configStorage,
		pipelineRunner: pipelineRunner,
	})
}
