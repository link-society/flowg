package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type DeleteRoleRequest struct {
	Role string `path:"role" minLength:"1"`
}

type DeleteRoleResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteRoleUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req DeleteRoleRequest,
				resp *DeleteRoleResponse,
			) error {
				err := ctrl.deps.AuthStorage.DeleteRole(ctx, req.Role)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to delete role",
						slog.String("role", req.Role),
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

	u.SetName("delete_role")
	u.SetTitle("Delete Role")
	u.SetDescription("Delete Role")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
