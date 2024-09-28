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

type DeletePipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type DeletePipelineResponse struct {
	Success bool `json:"success"`
}

func DeletePipelineUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req DeletePipelineRequest,
				resp *DeletePipelineResponse,
			) error {
				err := configStorage.DeletePipeline(ctx, req.Pipeline)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete pipeline",
						slog.String("channel", "api"),
						slog.String("pipeline", req.Pipeline),
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

	u.SetName("delete_pipeline")
	u.SetTitle("Delete Pipeline")
	u.SetDescription("Delete pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
