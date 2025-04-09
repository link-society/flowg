package consul

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	logger *slog.Logger
	opts   *ConsulServiceOptions
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult { return nil }

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult { return nil }

func (h *procHandler) Terminate(ctx actor.Context, err error) error { return nil }
