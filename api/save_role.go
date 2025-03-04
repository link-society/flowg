package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type SaveRoleRequest struct {
	Role   string   `path:"role" minLength:"1"`
	Scopes []string `json:"scopes"`
}

type SaveRoleResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) SaveRoleUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
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
						ctrl.logger.ErrorContext(
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

				err := ctrl.deps.AuthStorage.SaveRole(ctx, role)
				if err != nil {
					ctrl.logger.ErrorContext(
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
