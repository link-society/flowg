package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

// DeletePipelineDeps lists the dependencies of [NewDeletePipelineUsecase].
type DeletePipelineDeps struct {
	fx.In

	AuthStorage   authStorage.Storage
	ConfigStorage config.Storage
}

// DeletePipelineRequest identifies the pipeline to remove.
type DeletePipelineRequest struct {
	// Pipeline is the name of the pipeline to delete.
	Pipeline string `path:"pipeline" minLength:"1"`
}

// DeletePipelineResponse reports the outcome of the deletion.
type DeletePipelineResponse struct {
	// Success reports whether the pipeline was removed.
	Success bool `json:"success"`
}

// NewDeletePipelineUsecase removes a pipeline.
//
// Callers must have the write-pipelines permission. Deleting an absent pipeline is
// treated as a success.
func NewDeletePipelineUsecase(deps DeletePipelineDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req DeletePipelineRequest,
				resp *DeletePipelineResponse,
			) error {
				err := deps.ConfigStorage.DeletePipeline(ctx, req.Pipeline)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to delete pipeline",
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

func init() {
	routing.RegisterOperation(
		NewDeletePipelineUsecase,
		http.MethodDelete,
		"/api/v1/pipelines/{pipeline}",
	)
}
