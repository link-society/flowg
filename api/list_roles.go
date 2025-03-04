package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ListRolesRequest struct{}
type ListRolesResponse struct {
	Success bool          `json:"success"`
	Roles   []models.Role `json:"roles"`
}

func (ctrl *controller) ListRolesUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListRolesRequest,
				resp *ListRolesResponse,
			) error {
				roles, err := ctrl.deps.AuthStorage.ListRoles(ctx)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to list roles",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Roles = roles
				return nil
			},
		),
	)

	u.SetName("list_roles")
	u.SetTitle("List Roles")
	u.SetDescription("List known roles")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
