package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type SavePipelineRequest struct {
	Pipeline string             `path:"pipeline" minLength:"1"`
	Flow     models.FlowGraphV2 `json:"flow" required:"true"`
}

type SavePipelineResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) SavePipelineUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req SavePipelineRequest,
				resp *SavePipelineResponse,
			) error {
				if err := ctrl.deps.ConfigStorage.WritePipeline(ctx, req.Pipeline, &req.Flow); err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to save pipeline",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := ctrl.deps.PipelineRunner.InvalidateCachedBuild(ctx, req.Pipeline); err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to refresh pipeline cache after save",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := ctrl.deps.ClusterManager.BroadcastInvalidatePipelineCache(ctx, req.Pipeline); err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to broadcast pipeline cache invalidation after save",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("save_pipeline")
	u.SetTitle("Save Pipeline")
	u.SetDescription("Save pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
