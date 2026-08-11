package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ListAuthProvidersDeps lists the dependencies of [NewListAuthProvidersUsecase].
type ListAuthProvidersDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewListAuthProvidersUsecase enumerates all available auth providers.
func NewListAuthProvidersUsecase(deps ListAuthProvidersDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req schemas.ListAuthProvidersRequest,
			resp *schemas.ListAuthProvidersResponse,
		) error {
			authProviders, err := deps.AuthStorage.ListAuthProviders(ctx)
			if err != nil {
				logger.ErrorContext(
					ctx,
					"Failed to list auth providers",
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			authProviderInfos := make([]schemas.AuthProviderInfo, len(authProviders))	
			for i, ap := range authProviders {
				authProviderInfos[i] = schemas.AuthProviderInfo{
					Name:        ap.Name,
					DisplayName: ap.DisplayName,
				}
			}

			resp.Success = true
			resp.AuthProviders = authProviderInfos
			return nil
		},
	)

	u.SetName("list_auth_providers")
	u.SetTitle("List Auth Providers")
	u.SetDescription("List available auth providers")
	u.SetTags("auth")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListAuthProvidersUsecase,
		http.MethodGet,
		"/api/v1/auth-providers/",
		routing.Public(),
	)
}
