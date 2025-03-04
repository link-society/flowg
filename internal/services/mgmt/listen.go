package mgmt

import (
	"log/slog"

	"net"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
)

type listenerHandler struct {
	logger *slog.Logger

	bindAddress string

	listener net.Listener
}

func (h *listenerHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.logger.InfoContext(ctx, "Listen on Management interface")

	var err error
	h.listener, err = net.Listen("tcp", h.bindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to listen on Management interface",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *listenerHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *listenerHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
