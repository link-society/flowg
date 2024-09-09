package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type DeletePipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type DeletePipelineResponse struct {
	Success bool `json:"success"`
}

func DeletePipelineUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	pipelineSys := config.NewPipelineSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req DeletePipelineRequest,
				resp *DeletePipelineResponse,
			) error {
				err := pipelineSys.Delete(req.Pipeline)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete pipeline",
						"channel", "api",
						"pipeline", req.Pipeline,
						"error", err.Error(),
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
