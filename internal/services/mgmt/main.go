package mgmt

import (
	"log/slog"

	"crypto/tls"
	"net/url"

	"link-society.com/flowg/internal/utils/proctree"
)

type ServerOptions struct {
	BindAddress string
	TlsConfig   *tls.Config

	ClusterNodeID       string
	ClusterJoinNodeID   string
	ClusterJoinEndpoint *url.URL
}

func NewServer(opts *ServerOptions) proctree.Process {
	state := &state{
		logger: slog.Default().With(
			slog.String("channel", "mgmt"),
			slog.Group("mgmt",
				slog.String("bind", opts.BindAddress),
			),
		),

		opts: opts,
	}

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		proctree.NewProcess(&listenerHandler{state: state}),
		proctree.NewProcess(&clusterHandler{state: state}),
		proctree.NewProcess(&serverHandler{state: state}),
	)
}
