package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type DeleteUserRequest struct {
	User string `path:"user" minLength:"1"`
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

func DeleteUserUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req DeleteUserRequest,
				resp *DeleteUserResponse,
			) error {
				err := authDb.DeleteUser(req.User)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete user",
						"channel", "api",
						"user", req.User,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				return nil
			},
		),
	)

	u.SetName("delete_user")
	u.SetTitle("Delete User")
	u.SetDescription("Delete User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
