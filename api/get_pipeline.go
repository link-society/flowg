package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type GetPipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type GetPipelineResponse struct {
	Success bool              `json:"success"`
	Flow    *config.FlowGraph `json:"flow"`
}

func GetPipelineUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	pipelineSys := config.NewPipelineSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetPipelineRequest,
				resp *GetPipelineResponse,
			) error {
				flowGraph, err := pipelineSys.Parse(req.Pipeline)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get pipeline",
						"channel", "api",
						"pipeline", req.Pipeline,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Flow = flowGraph

				return nil
			},
		),
	)

	u.SetName("get_pipeline")
	u.SetTitle("Get Pipeline")
	u.SetDescription("Get pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
