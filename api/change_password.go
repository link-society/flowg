package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/app/logging"
	apiUtils "link-society.com/flowg/internal/utils/api"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" required:"true"`
	NewPassword string `json:"new_password" required:"true"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) ChangePasswordUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ChangePasswordRequest,
			resp *ChangePasswordResponse,
		) error {
			logging.MarkSensitive(ctx)

			user := apiUtils.GetContextUser(ctx)

			switch valid, err := ctrl.deps.AuthStorage.VerifyUserPassword(ctx, user.Name, req.OldPassword); {
			case err != nil:
				ctrl.logger.ErrorContext(
					ctx,
					"Failed to verify user password",
					slog.String("username", user.Name),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)

			case !valid:
				resp.Success = false
				return status.Wrap(errors.New("invalid credentials"), status.PermissionDenied)
			}

			if err := ctrl.deps.AuthStorage.SaveUser(ctx, *user, req.NewPassword); err != nil {
				ctrl.logger.ErrorContext(
					ctx,
					"Failed to update user password",
					slog.String("username", user.Name),
					slog.String("error", err.Error()),
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
