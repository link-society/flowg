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

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// SaveForwarderDeps lists the dependencies of [NewSaveForwarderUsecase].
type SaveForwarderDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	ConfigStorage  storage.ConfigStorage
	PipelineRunner pipelines.Runner
}

// NewSaveForwarderUsecase creates or overwrites a forwarder.
//
// Callers must have the write-forwarders permission. Persisting a forwarder
// invalidates cached pipeline builds so that subsequent runs use the new
// definition.
func NewSaveForwarderUsecase(deps SaveForwarderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_FORWARDERS,
			func(
				ctx context.Context,
				req schemas.SaveForwarderRequest,
				resp *schemas.SaveForwarderResponse,
			) error {
				err := deps.ConfigStorage.WriteForwarder(ctx, req.Forwarder, &req.Config)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := deps.PipelineRunner.InvalidateAllCachedBuilds(ctx); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to refresh pipeline cache after save",
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

	u.SetName("save_forwarder")
	u.SetTitle("Save Forwarder")
	u.SetDescription("Save forwarder")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSaveForwarderUsecase,
		http.MethodPut,
		"/api/v1/forwarders/{forwarder}",
	)
}
