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

// DeleteUserDeps lists the dependencies of [NewDeleteUserUsecase].
type DeleteUserDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewDeleteUserUsecase removes a user account.
//
// Callers must have the write-ACLs permission. Deleting an absent user is
// treated as a success.
func NewDeleteUserUsecase(deps DeleteUserDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req schemas.DeleteUserRequest,
				resp *schemas.DeleteUserResponse,
			) error {
				err := deps.AuthStorage.DeleteUser(ctx, req.User)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to delete user",
						slog.String("user", req.User),
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

	u.SetName("delete_user")
	u.SetTitle("Delete User")
	u.SetDescription("Delete User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewDeleteUserUsecase,
		http.MethodDelete,
		"/api/v1/users/{user}",
	)
}
