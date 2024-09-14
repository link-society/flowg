package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func LoginUsecase(authDb *auth.Database) usecase.Interactor {
	userSys := auth.NewUserSystem(authDb)

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req LoginRequest,
			resp *LoginResponse,
		) error {
			switch valid, err := userSys.VerifyUserPassword(req.Username, req.Password); {
			case err != nil:
				slog.ErrorContext(
					ctx,
					"Failed to verify user password",
					"channel", "api",
					"username", req.Username,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Unauthenticated)

			case !valid:
				resp.Success = false
				return status.Wrap(errors.New("invalid credentials"), status.Unauthenticated)

			case valid:
				token, err := auth.NewJWT(req.Username)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to create JWT",
						"channel", "api",
						"username", req.Username,
						"error", err.Error(),
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
