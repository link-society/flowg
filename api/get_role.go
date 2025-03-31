package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetRoleRequest struct {
	Role string `path:"role" minLength:"1"`
}

type GetRoleResponse struct {
	Success bool         `json:"success"`
	Role    *models.Role `json:"role"`
}

func (ctrl *controller) GetRoleUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req GetRoleRequest,
				resp *GetRoleResponse,
			) error {
				role, err := ctrl.deps.AuthStorage.FetchRole(ctx, req.Role)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get role",
						slog.String("role", req.Role),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Role = role

				return nil
			},
		),
	)

	u.SetName("get_role")
	u.SetTitle("Get Role")
	u.SetDescription("Get Role")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
