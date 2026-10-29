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

// WhoamiDeps lists the dependencies of [NewWhoamiUsecase].
type WhoamiDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage storage.AuthStorage
}

// NewWhoamiUsecase returns the profile and effective permissions of the caller.
//
// It lets a client discover who it is authenticated as and what it is allowed
// to do, typically to drive UI affordances. It requires authentication but no
// particular permission.
func NewWhoamiUsecase(deps WhoamiDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req schemas.WhoamiRequest,
			resp *schemas.WhoamiResponse,
		) error {
			resp.User = auth.GetContextUser(ctx)

			scopes, err := deps.AuthStorage.ListUserScopes(ctx, resp.User.Name)
			if err != nil {
				deps.Logger.ErrorContext(
					ctx,
					"Failed to fetch user scopes",
					slog.String("channel", "api"),
					slog.String("error", err.Error()),
				)
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Permissions = models.PermissionsFromScopes(scopes)
			return nil
		},
	)

	u.SetName("whoami")
	u.SetTitle("Fetch current profile")
	u.SetDescription("Fetch the profile of the currently authenticated user")
	u.SetTags("auth")

	u.SetExpectedErrors(status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewWhoamiUsecase,
		http.MethodGet,
		"/api/v1/auth/whoami",
	)
}
