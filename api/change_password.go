package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

func ChangePasswordUsecase(authDb *auth.Database) usecase.Interactor {
	userSys := auth.NewUserSystem(authDb)

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ChangePasswordRequest,
			resp *ChangePasswordResponse,
		) error {
			user := auth.GetContextUser(ctx)

			switch valid, err := userSys.VerifyUserPassword(user.Name, req.OldPassword); {
			case err != nil:
				slog.ErrorContext(
					ctx,
					"Failed to verify user password",
					"channel", "api",
					"username", user.Name,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)

			case !valid:
				resp.Success = false
				return status.Wrap(errors.New("invalid credentials"), status.PermissionDenied)
			}

			if err := userSys.SaveUser(*user, req.NewPassword); err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to update user password",
					"channel", "api",
					"username", user.Name,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true

			return nil
		},
	)

	u.SetName("change_password")
	u.SetTitle("Change Password")
	u.SetDescription("Change Password of current user")
	u.SetTags("auth")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Internal)

	return u
}
