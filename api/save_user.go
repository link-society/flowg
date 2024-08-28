package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type SaveUserRequest struct {
	User     string   `path:"user" minLength:"1"`
	Roles    []string `json:"roles"`
	Password string   `json:"password"`
}

type SaveUserResponse struct {
	Success bool `json:"success"`
}

func SaveUserUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req SaveUserRequest,
				resp *SaveUserResponse,
			) error {
				user := auth.User{
					Name:  req.User,
					Roles: req.Roles,
				}

				err := authDb.SaveUser(user, req.Password)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save user",
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

	u.SetName("save_user")
	u.SetTitle("Save User")
	u.SetDescription("Save User")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
