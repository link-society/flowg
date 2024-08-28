package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type SaveRoleRequest struct {
	Role   string   `path:"role" minLength:"1"`
	Scopes []string `json:"scopes"`
}

type SaveRoleResponse struct {
	Success bool `json:"success"`
}

func SaveRoleUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req SaveRoleRequest,
				resp *SaveRoleResponse,
			) error {
				scopes := make([]auth.Scope, 0, len(req.Scopes))

				for i, scopeName := range req.Scopes {
					scope, err := auth.ParseScope(scopeName)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to parse scope",
							"channel", "api",
							"scope", scopeName,
							"error", err.Error(),
						)

						resp.Success = false
						return status.Wrap(err, status.InvalidArgument)
					}

					scopes[i] = scope
				}

				role := auth.Role{
					Name:   req.Role,
					Scopes: scopes,
				}

				err := authDb.SaveRole(role)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save role",
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

	u.SetName("save_role")
	u.SetTitle("Save Role")
	u.SetDescription("Save Role")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}
