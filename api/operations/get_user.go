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

// GetUserDeps lists the dependencies of [NewGetUserUsecase].
type GetUserDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewGetUserUsecase returns a single user account along with its roles.
//
// Callers must have the read-ACLs permission. Requesting an unknown user yields
// a not-found error.
func NewGetUserUsecase(deps GetUserDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req schemas.GetUserRequest,
				resp *schemas.GetUserResponse,
			) error {
				user, err := deps.AuthStorage.FetchUser(ctx, req.Username)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get user",
						slog.String("user", req.Username),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.User = user

				return nil
			},
		),
	)

	u.SetName("get_user")
	u.SetTitle("Get User")
	u.SetDescription("Get User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetUserUsecase,
		http.MethodGet,
		"/api/v1/users/{user}",
	)
}
