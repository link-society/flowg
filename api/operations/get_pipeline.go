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

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// GetPipelineDeps lists the dependencies of [NewGetPipelineUsecase].
type GetPipelineDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// GetPipelineRequest identifies the pipeline to retrieve.
type GetPipelineRequest struct {
	// Pipeline is the name of the pipeline to read.
	Pipeline string `path:"pipeline" minLength:"1"`
}

// GetPipelineResponse carries the flow graph of the requested pipeline.
type GetPipelineResponse struct {
	// Success reports whether the pipeline was found and returned.
	Success bool `json:"success"`
	// Flow is the pipeline's flow graph definition.
	Flow *models.FlowGraphV2 `json:"flow"`
}

// NewGetPipelineUsecase returns the flow graph of a single pipeline.
//
// Callers must have the read-pipelines permission. Requesting an unknown pipeline
// yields a not-found error.
func NewGetPipelineUsecase(deps GetPipelineDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetPipelineRequest,
				resp *GetPipelineResponse,
			) error {
				flowGraph, err := deps.ConfigStorage.ReadPipeline(ctx, req.Pipeline)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewGetPipelineUsecase,
		http.MethodGet,
		"/api/v1/pipelines/{pipeline}",
	)
}
