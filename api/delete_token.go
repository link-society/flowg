package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"
)

type DeleteTokenRequest struct {
	TokenUUID string `path:"token-uuid" format:"uuid"`
}

type DeleteTokenResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteTokenUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req DeleteTokenRequest,
			resp *DeleteTokenResponse,
		) error {
			user := apiUtils.GetContextUser(ctx)

			err := ctrl.deps.AuthStorage.DeleteToken(ctx, user.Name, req.TokenUUID)
			if err != nil {
				ctrl.logger.ErrorContext(
					ctx,
					"Failed to delete token",
					slog.String("user", user.Name),
					slog.String("token-uuid", req.TokenUUID),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true

			return nil
		},
	)

	u.SetName("delete_token")
	u.SetTitle("Delete Token")
	u.SetDescription("Delete Personal Access Token UUIDs for the current user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.Internal)

	return u
}
