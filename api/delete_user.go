package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type DeleteUserRequest struct {
	User string `path:"user" minLength:"1"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteUserUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req DeleteUserRequest,
				resp *DeleteUserResponse,
			) error {
				err := ctrl.deps.AuthStorage.DeleteUser(ctx, req.User)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to delete user",
						slog.String("user", req.User),
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

	u.SetName("delete_user")
	u.SetTitle("Delete User")
	u.SetDescription("Delete User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
