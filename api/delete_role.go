package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type DeleteRoleRequest struct {
	Role string `path:"role" minLength:"1"`
}

type DeleteRoleResponse struct {
	Success bool `json:"success"`
}

func DeleteRoleUsecase(authDb *auth.Database) usecase.Interactor {
	roleSys := auth.NewRoleSystem(authDb)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req DeleteRoleRequest,
				resp *DeleteRoleResponse,
			) error {
				err := roleSys.DeleteRole(req.Role)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete role",
						"channel", "api",
						"role", req.Role,
						"error", err.Error(),
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
