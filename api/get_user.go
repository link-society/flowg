package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type GetUserRequest struct {
	Username string `path:"user" minLength:"1"`
}

type GetUserResponse struct {
	Success bool         `json:"success"`
	User    *models.User `json:"user"`
}

func (ctrl *controller) GetUserUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req GetUserRequest,
				resp *GetUserResponse,
			) error {
				user, err := ctrl.deps.AuthStorage.FetchUser(ctx, req.Username)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get user",
						slog.String("user", req.Username),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.User = user

				return nil
			},
		),
	)

	u.SetName("get_user")
	u.SetTitle("Get User")
	u.SetDescription("Get User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
