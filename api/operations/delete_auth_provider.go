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

// DeleteAuthProviderDeps lists the dependencies of [NewDeleteAuthProviderUsecase].
type DeleteAuthProviderDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewDeleteAuthProviderUsecase removes an auth provider.
//
// Callers must have the write-auth-providers permission. Deleting an absent auth provider is
// treated as a success.
func NewDeleteAuthProviderUsecase(deps DeleteAuthProviderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_AUTH_PROVIDERS,
			func(
				ctx context.Context,
				req schemas.DeleteAuthProviderRequest,
				resp *schemas.DeleteAuthProviderResponse,
			) error {
				err := deps.AuthStorage.DeleteAuthProvider(ctx, req.AuthProvider)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to delete auth provider",
						slog.String("auth_provider", req.AuthProvider),
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

	u.SetName("delete_auth_provider")
	u.SetTitle("Delete Auth Provider")
	u.SetDescription("Delete Auth Provider")
	u.SetTags("auth")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewDeleteAuthProviderUsecase,
		http.MethodDelete,
		"/api/v1/auth-providers/{auth_provider}",
	)
}
