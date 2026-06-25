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

// ListRolesDeps lists the dependencies of [NewListRolesUsecase].
type ListRolesDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// ListRolesRequest is empty: listing roles takes no parameters.
type ListRolesRequest struct{}

// ListRolesResponse carries every known role with its permissions.
type ListRolesResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Roles holds every configured role and its granted permissions.
	Roles []models.Role `json:"roles"`
}

// NewListRolesUsecase enumerates all configured roles with their permissions.
//
// Callers must have the read-ACLs permission.
func NewListRolesUsecase(deps ListRolesDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListRolesRequest,
				resp *ListRolesResponse,
			) error {
				roles, err := deps.AuthStorage.ListRoles(ctx)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewListRolesUsecase,
		http.MethodGet,
		"/api/v1/roles",
	)
}
