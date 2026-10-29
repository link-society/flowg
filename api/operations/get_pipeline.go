package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	"go.link-society.com/flowg/internal/models"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// GetPipelineDeps lists the dependencies of [NewGetPipelineUsecase].
type GetPipelineDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewGetPipelineUsecase returns the flow graph of a single pipeline.
//
// Callers must have the read-pipelines permission. Requesting an unknown pipeline
// yields a not-found error.
func NewGetPipelineUsecase(deps GetPipelineDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req schemas.GetPipelineRequest,
				resp *schemas.GetPipelineResponse,
			) error {
				flowGraph, err := deps.ConfigStorage.ReadPipeline(ctx, req.Pipeline)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to get pipeline",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}
				if flowGraph == nil {
					resp.Success = false
					return status.NotFound
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

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetPipelineUsecase,
		http.MethodGet,
		"/api/v1/pipelines/{pipeline}",
	)
}
