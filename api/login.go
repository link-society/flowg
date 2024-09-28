package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	authUtils "link-society.com/flowg/internal/utils/auth"

	"link-society.com/flowg/internal/storage/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func LoginUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req LoginRequest,
			resp *LoginResponse,
		) error {
			switch valid, err := authStorage.VerifyUserPassword(ctx, req.Username, req.Password); {
			case err != nil:
				slog.ErrorContext(
					ctx,
					"Failed to verify user password",
					slog.String("channel", "api"),
					slog.String("username", req.Username),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Unauthenticated)

			case !valid:
				resp.Success = false
				return status.Wrap(errors.New("invalid credentials"), status.Unauthenticated)

			case valid:
				token, err := authUtils.NewJWT(req.Username)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to create JWT",
						slog.String("channel", "api"),
						slog.String("username", req.Username),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Unauthenticated)
				}

				resp.Success = true
				resp.Token = token
			}

			return nil
		},
	)

	u.SetName("login")
	u.SetTitle("Authenticate")
	u.SetDescription("Create new Session cookie")
	u.SetTags("auth")

	u.SetExpectedErrors(status.Unauthenticated)

	return u
}
