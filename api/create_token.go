package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type CreateTokenRequest struct {
	User string `path:"user" minLength:"1"`
}

type CreateTokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func CreateTokenUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req CreateTokenRequest,
				resp *CreateTokenResponse,
			) error {
				user, err := authDb.GetUser(req.User)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get user",
						"channel", "api",
						"user", req.User,
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
						"user", req.User,
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
						"user", req.User,
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
						"user", req.User,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Token = token

				return nil
			},
		),
	)

	u.SetName("create_token")
	u.SetTitle("Create Token")
	u.SetDescription("Create a new Personal Access Token for a user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
