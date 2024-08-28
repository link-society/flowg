package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type CreateTokenRequest struct{}

type CreateTokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func CreateTokenUsecase(authDb *auth.Database) usecase.Interactor {
	tokenSys := auth.NewTokenSystem(authDb)

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req CreateTokenRequest,
			resp *CreateTokenResponse,
		) error {
			user := auth.GetContextUser(ctx)

			token, err := tokenSys.CreateToken(user.Name)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to create token",
					"channel", "api",
					"user", user.Name,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Token = token

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
