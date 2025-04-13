package server

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"

	"link-society.com/flowg/internal/app/bootstrap"
)

type bootstrapProcHandler struct {
	logger *slog.Logger

	authStorage   *auth.Storage
	configStorage *config.Storage
}

var _ proctree.ProcessHandler = (*bootstrapProcHandler)(nil)

func (h *bootstrapProcHandler) Init(ctx actor.Context) proctree.ProcessResult {
	err := bootstrap.DefaultRolesAndUsers(ctx, h.authStorage)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to bootstrap default roles and users",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	err = bootstrap.DefaultPipeline(ctx, h.configStorage)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to bootstrap default pipeline",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *bootstrapProcHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *bootstrapProcHandler) Terminate(ctx actor.Context, err error) error {
	return err
}
