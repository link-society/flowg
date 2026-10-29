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

	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/pkg/models"
)

// GetForwarderDeps lists the dependencies of [NewGetForwarderUsecase].
type GetForwarderDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewGetForwarderUsecase returns the definition of a single forwarder.
//
// Callers must have the read-forwarders permission. Requesting an unknown
// forwarder yields a not-found error.
func NewGetForwarderUsecase(deps GetForwarderDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_FORWARDERS,
			func(
				ctx context.Context,
				req schemas.GetForwarderRequest,
				resp *schemas.GetForwarderResponse,
			) error {
				forwarder, err := deps.ConfigStorage.ReadForwarder(ctx, req.Forwarder)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to get forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}
				if forwarder == nil {
					resp.Success = false
					return status.NotFound
				}

				resp.Success = true
				resp.Forwarder = forwarder

				return nil
			},
		),
	)

	u.SetName("get_forwarder")
	u.SetTitle("Get Forwarder")
	u.SetDescription("Get forwarder")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetForwarderUsecase,
		http.MethodGet,
		"/api/v1/forwarders/{forwarder}",
	)
}
