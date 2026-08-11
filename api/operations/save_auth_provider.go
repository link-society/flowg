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

// SaveAuthProviderDeps lists the dependencies of [NewSaveAuthProviderUsecase].
type SaveAuthProviderDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewSaveAuthProviderUsecase creates or overwrites an auth provider.
//
// Callers must have the write-auth-providers permission.
func NewSaveAuthProviderUsecase(deps SaveAuthProviderDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_AUTH_PROVIDERS,
			func(
				ctx context.Context,
				req schemas.SaveAuthProviderRequest,
				resp *schemas.SaveAuthProviderResponse,
			) error {
				req.Config.Name = req.AuthProvider
				err := deps.AuthStorage.SaveAuthProvider(ctx, req.Config)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save auth provider",
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

	u.SetName("save_auth_provider")
	u.SetTitle("Save Auth Provider")
	u.SetDescription("Save auth provider")
	u.SetTags("auth")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSaveAuthProviderUsecase,
		http.MethodPut,
		"/api/v1/auth-providers/{auth_provider}",
	)
}
