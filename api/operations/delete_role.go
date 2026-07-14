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

// DeleteRoleDeps lists the dependencies of [NewDeleteRoleUsecase].
type DeleteRoleDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewDeleteRoleUsecase removes a role.
//
// Callers must have the write-ACLs permission. Deleting an absent role is
// treated as a success.
func NewDeleteRoleUsecase(deps DeleteRoleDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req schemas.DeleteRoleRequest,
				resp *schemas.DeleteRoleResponse,
			) error {
				err := deps.AuthStorage.DeleteRole(ctx, req.Role)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewDeleteRoleUsecase,
		http.MethodDelete,
		"/api/v1/roles/{role}",
	)
}
