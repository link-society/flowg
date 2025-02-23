package mgmt

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"
)

func NewServer(bindAddress string, tlsConfig *tls.Config) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "mgmt"),
			slog.Group("mgmt",
				slog.String("bind", bindAddress),
			),
		),

		bindAddress: bindAddress,
		tlsConfig:   tlsConfig,
	})
}
