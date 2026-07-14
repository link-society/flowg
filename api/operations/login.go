package operations

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// LoginDeps lists the dependencies of [NewLoginUsecase].
type LoginDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewLoginUsecase authenticates a user by password and issues a session token.
//
// It is the only operation that requires no prior authentication. Invalid
// credentials are reported as an unauthenticated error, never revealing whether
// the username exists.
func NewLoginUsecase(deps LoginDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req schemas.LoginRequest,
			resp *schemas.LoginResponse,
		) error {
			switch valid, err := deps.AuthStorage.VerifyUserPassword(ctx, req.Username, req.Password); {
			case err != nil:
				logger.ErrorContext(
					ctx,
					"Failed to verify user password",
					slog.String("username", req.Username),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Unauthenticated)

			case !valid:
				resp.Success = false
				return status.Wrap(errors.New("invalid credentials"), status.Unauthenticated)

			case valid:
				token, err := auth.NewJWT(req.Username)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to create JWT",
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

func init() {
	routing.RegisterOperation(
		NewLoginUsecase,
		http.MethodPost,
		"/api/v1/auth/login",
		routing.Public(),
	)
}
