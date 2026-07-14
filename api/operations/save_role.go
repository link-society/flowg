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

// SaveRoleDeps lists the dependencies of [NewSaveRoleUsecase].
type SaveRoleDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewSaveRoleUsecase creates or overwrites a role with a given set of permissions.
//
// Callers must have the write-ACLs permission. An unknown permission name in
// the request is reported as an invalid-argument error.
func NewSaveRoleUsecase(deps SaveRoleDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req schemas.SaveRoleRequest,
				resp *schemas.SaveRoleResponse,
			) error {
				scopes := make([]models.Scope, len(req.Scopes))

				for i, scopeName := range req.Scopes {
					scope, err := models.ParseScope(scopeName)
					if err != nil {
						logger.ErrorContext(
							ctx,
							"Failed to parse scope",
							slog.String("scope", scopeName),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						return status.Wrap(err, status.InvalidArgument)
					}

					scopes[i] = scope
				}

				role := models.Role{
					Name:   req.Role,
					Scopes: scopes,
				}

				err := deps.AuthStorage.SaveRole(ctx, role)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save role",
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

	u.SetName("save_role")
	u.SetTitle("Save Role")
	u.SetDescription("Save Role")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSaveRoleUsecase,
		http.MethodPut,
		"/api/v1/roles/{role}",
	)
}
