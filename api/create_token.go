package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"
)

type CreateTokenRequest struct{}

type CreateTokenResponse struct {
	Success   bool   `json:"success"`
	Token     string `json:"token"`
	TokenUUID string `json:"token_uuid"`
}

func (ctrl *controller) CreateTokenUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req CreateTokenRequest,
			resp *CreateTokenResponse,
		) error {
			user := apiUtils.GetContextUser(ctx)

			token, tokenUuid, err := ctrl.deps.AuthStorage.CreateToken(ctx, user.Name)
			if err != nil {
				ctrl.logger.ErrorContext(
					ctx,
					"Failed to create token",
					slog.String("user", user.Name),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Token = token
			resp.TokenUUID = tokenUuid

			return nil
		},
	)

	u.SetName("create_token")
	u.SetTitle("Create Token")
	u.SetDescription("Create a new Personal Access Token for the current user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.NotFound, status.Internal)

	return u
}
