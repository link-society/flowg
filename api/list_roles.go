package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/auth"
)

type ListRolesRequest struct{}
type ListRolesResponse struct {
	Success bool        `json:"success"`
	Roles   []auth.Role `json:"roles"`
}

func ListRolesUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListRolesRequest,
				resp *ListRolesResponse,
			) error {
				roleNames, err := authDb.ListRoles()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list roles",
						"channel", "api",
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				roles := make([]auth.Role, 0, len(roleNames))

				for i, roleName := range roleNames {
					role, err := authDb.GetRole(roleName)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to get role",
							"channel", "api",
							"role", roleName,
							"error", err.Error(),
						)

						resp.Success = false
						return status.Wrap(err, status.Internal)
					}

					roles[i] = role
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
