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

// PatchUserRolesDeps lists the dependencies of [NewPatchUserRolesUsecase].
type PatchUserRolesDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// PatchUserRolesRequest carries the user and the roles to assign it.
type PatchUserRolesRequest struct {
	// User is the name of the account to update.
	User string `path:"user" minLength:"1"`
	// Roles are the names of the roles to assign, replacing the current set.
	Roles []string `json:"roles" required:"true" items.minLength:"1"`
}

// PatchUserRolesResponse reports the outcome of the update.
type PatchUserRolesResponse struct {
	// Success reports whether the roles were updated.
	Success bool `json:"success"`
}

// NewPatchUserRolesUsecase replaces the roles assigned to a user without touching
// its password.
//
// It is the role-management counterpart to [NewSaveUserUsecase], used
// when the caller does not know or want to reset the password. Callers must
// have the write-ACLs permission.
func NewPatchUserRolesUsecase(deps PatchUserRolesDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req PatchUserRolesRequest,
				resp *PatchUserRolesResponse,
			) error {
				user := models.User{
					Name:  req.User,
					Roles: req.Roles,
				}

				err := deps.AuthStorage.PatchUserRoles(ctx, user)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to patch user roles",
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

	u.SetName("patch_user_roles")
	u.SetTitle("Patch User Roles")
	u.SetDescription("Patch User Roles")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewPatchUserRolesUsecase,
		http.MethodPatch,
		"/api/v1/users/{user}",
	)
}
