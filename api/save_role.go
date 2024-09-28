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

type SaveRoleRequest struct {
	Role   string   `path:"role" minLength:"1"`
	Scopes []string `json:"scopes"`
}

type SaveRoleResponse struct {
	Success bool `json:"success"`
}

func SaveRoleUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req SaveRoleRequest,
				resp *SaveRoleResponse,
			) error {
				scopes := make([]models.Scope, len(req.Scopes))

				for i, scopeName := range req.Scopes {
					scope, err := models.ParseScope(scopeName)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to parse scope",
							slog.String("channel", "api"),
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

				err := authStorage.SaveRole(ctx, role)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save role",
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

	u.SetName("save_role")
	u.SetTitle("Save Role")
	u.SetDescription("Save Role")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}
