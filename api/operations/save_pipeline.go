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
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// SavePipelineDeps lists the dependencies of [NewSavePipelineUsecase].
type SavePipelineDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	ConfigStorage  storage.ConfigStorage
	PipelineRunner pipelines.Runner
}

// SavePipelineRequest carries the pipeline name and its new flow graph.
type SavePipelineRequest struct {
	// Pipeline is the name of the pipeline to create or overwrite.
	Pipeline string `path:"pipeline" minLength:"1"`
	// Flow is the flow graph to store under that name.
	Flow models.FlowGraphV2 `json:"flow" required:"true"`
}

// SavePipelineResponse reports the outcome of the save.
type SavePipelineResponse struct {
	// Success reports whether the pipeline was persisted.
	Success bool `json:"success"`
}

// NewSavePipelineUsecase creates or overwrites a pipeline.
//
// Callers must have the write-pipelines permission. Persisting a pipeline
// invalidates its cached build so that subsequent runs use the new flow graph.
func NewSavePipelineUsecase(deps SavePipelineDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req SavePipelineRequest,
				resp *SavePipelineResponse,
			) error {
				if err := deps.ConfigStorage.WritePipeline(ctx, req.Pipeline, &req.Flow); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save pipeline",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := deps.PipelineRunner.InvalidateCachedBuild(ctx, req.Pipeline); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to refresh pipeline cache after save",
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

	u.SetName("save_pipeline")
	u.SetTitle("Save Pipeline")
	u.SetDescription("Save pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSavePipelineUsecase,
		http.MethodPut,
		"/api/v1/pipelines/{pipeline}",
	)
}
