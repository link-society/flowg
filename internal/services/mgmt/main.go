package mgmt

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config
}

func NewServer(opts *ServerOptions) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "mgmt"),
			slog.Group("mgmt",
				slog.String("bind", opts.BindAddress),
			),
		),

		opts: opts,
	})
}
