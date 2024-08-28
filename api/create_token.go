package api

import (
	"context"
	"errors"
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
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req CreateTokenRequest,
			resp *CreateTokenResponse,
		) error {
			username := auth.GetContextUser(ctx)
			user, err := authDb.GetUser(username)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to get user",
					"channel", "api",
					"user", username,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			if user == nil {
				slog.ErrorContext(
					ctx,
					"User not found",
					"channel", "api",
					"user", username,
				)

				resp.Success = false
				return status.Wrap(errors.New("user not found"), status.NotFound)
			}

			token, err := auth.NewToken(32)
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

			err = authDb.AddPersonalAccessToken(user.Name, token)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to save token",
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

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
