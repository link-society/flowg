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

	authStorage "link-society.com/flowg/internal/storage/auth"
)

// SaveUserDeps lists the dependencies of [NewSaveUserUsecase].
type SaveUserDeps struct {
	fx.In

	AuthStorage authStorage.Storage
}

// SaveUserRequest carries a user account and its initial password.
type SaveUserRequest struct {
	// User is the name of the account to create or overwrite.
	User string `path:"user" minLength:"1"`
	// Roles are the names of the roles to assign the user.
	Roles []string `json:"roles" required:"true"`
	// Password is the account's password; it is stored hashed.
	Password string `json:"password" required:"true"`
}

// SaveUserResponse reports the outcome of the save.
type SaveUserResponse struct {
	// Success reports whether the user was persisted.
	Success bool `json:"success"`
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
				req SaveUserRequest,
				resp *SaveUserResponse,
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
