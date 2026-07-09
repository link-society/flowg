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

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// GetRoleDeps lists the dependencies of [NewGetRoleUsecase].
type GetRoleDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// GetRoleRequest identifies the role to retrieve.
type GetRoleRequest struct {
	// Role is the name of the role to read.
	Role string `path:"role" minLength:"1"`
}

// GetRoleResponse carries the definition of the requested role.
type GetRoleResponse struct {
	// Success reports whether the role was found and returned.
	Success bool `json:"success"`
	// Role is the role and its granted permissions.
	Role *models.Role `json:"role"`
}

// NewGetRoleUsecase returns the definition of a single role.
//
// Callers must have the read-ACLs permission. Requesting an unknown role yields
// a not-found error.
func NewGetRoleUsecase(deps GetRoleDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req GetRoleRequest,
				resp *GetRoleResponse,
			) error {
				role, err := deps.AuthStorage.FetchRole(ctx, req.Role)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewGetRoleUsecase,
		http.MethodGet,
		"/api/v1/roles/{role}",
	)
}
