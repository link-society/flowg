package consul

import (
	"log/slog"

	"link-society.com/flowg/internal/utils/proctree"
)

type ConsulServiceOptions struct {
	NodeId      string
	NodeAddress string
	NodePort    string
	ServiceName string
	ConsulUrl   string
}

func NewConsulService(opts *ConsulServiceOptions) proctree.Process {
	return proctree.NewProcess(&procHandler{
		logger: slog.Default().With(
			slog.String("channel", "consul"),
			slog.Group("consul",
				slog.String("consulUrl", opts.ConsulUrl),
			),
		),

		opts: opts,
	})
}
