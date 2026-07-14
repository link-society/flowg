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
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ListPipelinesDeps lists the dependencies of [NewListPipelinesUsecase].
type ListPipelinesDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewListPipelinesUsecase enumerates the names of all configured pipelines.
//
// Callers must have the read-pipelines permission.
func NewListPipelinesUsecase(deps ListPipelinesDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req schemas.ListPipelinesRequest,
				resp *schemas.ListPipelinesResponse,
			) error {
				pipelines, err := deps.ConfigStorage.ListPipelines(ctx)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to list pipelines",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Pipelines = pipelines

				return nil
			},
		),
	)

	u.SetName("list_pipelines")
	u.SetTitle("List Pipelines")
	u.SetDescription("List pipelines")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListPipelinesUsecase,
		http.MethodGet,
		"/api/v1/pipelines",
	)
}
