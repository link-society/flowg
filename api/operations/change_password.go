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
	applog "link-society.com/flowg/internal/app/logging"

	"link-society.com/flowg/internal/storage"
)

// ChangePasswordDeps lists the dependencies of [NewChangePasswordUsecase].
type ChangePasswordDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// ChangePasswordRequest carries the caller's current and desired passwords.
type ChangePasswordRequest struct {
	// OldPassword is the caller's current password, required to authorize the
	// change.
	OldPassword string `json:"old_password" required:"true" minLength:"1"`
	// NewPassword is the password to set.
	NewPassword string `json:"new_password" required:"true" minLength:"1"`
}

// ChangePasswordResponse reports the outcome of the change.
type ChangePasswordResponse struct {
	// Success reports whether the password was updated.
	Success bool `json:"success"`
}

// NewChangePasswordUsecase lets the authenticated user change their own password.
//
// It re-verifies the current password before applying the change, so a stolen
// session alone cannot lock the owner out. It requires authentication but no
// particular permission; a wrong current password is reported as
// permission-denied. The request is marked sensitive so the passwords are kept
// out of FlowG's own logs.
func NewChangePasswordUsecase(deps ChangePasswordDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ChangePasswordRequest,
			resp *ChangePasswordResponse,
		) error {
			applog.MarkSensitive(ctx)

			user := auth.GetContextUser(ctx)

			switch valid, err := deps.AuthStorage.VerifyUserPassword(ctx, user.Name, req.OldPassword); {
			case err != nil:
				logger.ErrorContext(
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

			if err := deps.AuthStorage.SaveUser(ctx, *user, req.NewPassword); err != nil {
				logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewChangePasswordUsecase,
		http.MethodPost,
		"/api/v1/auth/change-password",
	)
}
