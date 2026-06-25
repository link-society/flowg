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

	"link-society.com/flowg/internal/storage"
)

// ListForwardersDeps lists the dependencies of [NewListForwardersUsecase].
type ListForwardersDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// ListForwardersRequest is empty: listing forwarders takes no parameters.
type ListForwardersRequest struct{}

// ListForwardersResponse carries the names of the available forwarders.
type ListForwardersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Forwarders holds the name of every configured forwarder.
	Forwarders []string `json:"forwarders"`
}

// NewListForwardersUsecase enumerates the names of all configured forwarders.
//
// Callers must have the read-forwarders permission.
func NewListForwardersUsecase(deps ListForwardersDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_FORWARDERS,
			func(
				ctx context.Context,
				req ListForwardersRequest,
				resp *ListForwardersResponse,
			) error {
				forwarders, err := deps.ConfigStorage.ListForwarders(ctx)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to list forwarders",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Forwarders = forwarders

				return nil
			},
		),
	)

	u.SetName("list_forwarders")
	u.SetTitle("List Forwarders")
	u.SetDescription("List forwarders")
	u.SetTags("forwarders")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListForwardersUsecase,
		http.MethodGet,
		"/api/v1/forwarders",
	)
}
