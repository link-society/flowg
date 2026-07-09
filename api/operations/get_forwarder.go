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

// GetForwarderDeps lists the dependencies of [NewGetForwarderUsecase].
type GetForwarderDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// GetForwarderRequest identifies the forwarder to retrieve.
type GetForwarderRequest struct {
	// Forwarder is the name of the forwarder to read.
	Forwarder string `path:"forwarder" minLength:"1"`
}

// GetForwarderResponse carries the definition of the requested forwarder.
type GetForwarderResponse struct {
	// Success reports whether the forwarder was found and returned.
	Success bool `json:"success"`
	// Forwarder is the forwarder's configuration.
	Forwarder *models.ForwarderV2 `json:"forwarder"`
}

// NewGetForwarderUsecase returns the definition of a single forwarder.
//
// Callers must have the read-forwarders permission. Requesting an unknown
// forwarder yields a not-found error.
func NewGetForwarderUsecase(deps GetForwarderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_FORWARDERS,
			func(
				ctx context.Context,
				req GetForwarderRequest,
				resp *GetForwarderResponse,
			) error {
				forwarder, err := deps.ConfigStorage.ReadForwarder(ctx, req.Forwarder)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get forwarder",
						slog.String("forwarder", req.Forwarder),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
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

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetForwarderUsecase,
		http.MethodGet,
		"/api/v1/forwarders/{forwarder}",
	)
}
