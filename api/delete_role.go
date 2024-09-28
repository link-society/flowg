package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
)

type DeleteRoleRequest struct {
	Role string `path:"role" minLength:"1"`
}

type DeleteRoleResponse struct {
	Success bool `json:"success"`
}

func DeleteRoleUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req DeleteRoleRequest,
				resp *DeleteRoleResponse,
			) error {
				err := authStorage.DeleteRole(ctx, req.Role)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete role",
						slog.String("channel", "api"),
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
