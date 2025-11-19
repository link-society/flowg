package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type PatchUserRolesRequest struct {
	User  string   `path:"user" minLength:"1"`
	Roles []string `json:"roles" required:"true" items.minLength:"1"`
}

type PatchUserRolesResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) PatchUserRolesUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
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

				err := ctrl.deps.AuthStorage.PatchUserRoles(ctx, user)
				if err != nil {
					ctrl.logger.ErrorContext(
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
