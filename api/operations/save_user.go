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

// SaveUserDeps lists the dependencies of [NewSaveUserUsecase].
type SaveUserDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewSaveUserUsecase creates or overwrites a user account, including its password.
//
// Callers must have the write-ACLs permission. To change roles without
// resetting the password, use [NewPatchUserRolesUsecase] instead.
func NewSaveUserUsecase(deps SaveUserDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req schemas.SaveUserRequest,
				resp *schemas.SaveUserResponse,
			) error {
				user := models.User{
					Name:  req.User,
					Roles: req.Roles,
				}

				err := deps.AuthStorage.SaveUser(ctx, user, req.Password)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save user",
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

	u.SetName("save_user")
	u.SetTitle("Save User")
	u.SetDescription("Save User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSaveUserUsecase,
		http.MethodPut,
		"/api/v1/users/{user}",
	)
}
