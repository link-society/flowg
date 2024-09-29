package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

type GetPipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type GetPipelineResponse struct {
	Success bool                `json:"success"`
	Flow    *models.FlowGraphV1 `json:"flow"`
}

func GetPipelineUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetPipelineRequest,
				resp *GetPipelineResponse,
			) error {
				flowGraph, err := configStorage.ReadPipeline(ctx, req.Pipeline)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get pipeline",
						slog.String("channel", "api"),
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
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
