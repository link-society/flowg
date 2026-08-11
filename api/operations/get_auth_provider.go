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

// GetAuthProviderDeps lists the dependencies of [NewGetAuthProviderUsecase].
type GetAuthProviderDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewGetAuthProviderUsecase returns a single auth provider along with its details.
//
// Callers must have the read-auth-providers permission. Requesting an unknown auth provider yields
// a not-found error.
func NewGetAuthProviderUsecase(deps GetAuthProviderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_AUTH_PROVIDERS,
			func(
				ctx context.Context,
				req schemas.GetAuthProvidersRequest,
				resp *schemas.GetAuthProvidersResponse,
			) error {
				authProvider, err := deps.AuthStorage.FetchAuthProvider(ctx, req.AuthProvider)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get auth provider",
						slog.String("auth_provider", req.AuthProvider),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.AuthProvider = authProvider

				return nil
			},
		),
	)

	u.SetName("get_auth_provider")
	u.SetTitle("Get Auth Provider")
	u.SetDescription("Get Auth Provider")
	u.SetTags("auth")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetAuthProviderUsecase,
		http.MethodGet,
		"/api/v1/auth-providers/{auth_provider}",
	)
}
