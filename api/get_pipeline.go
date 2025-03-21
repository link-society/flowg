package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetPipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type GetPipelineResponse struct {
	Success bool                `json:"success"`
	Flow    *models.FlowGraphV2 `json:"flow"`
}

func (ctrl *controller) GetPipelineUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetPipelineRequest,
				resp *GetPipelineResponse,
			) error {
				flowGraph, err := ctrl.deps.ConfigStorage.ReadPipeline(ctx, req.Pipeline)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get pipeline",
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
