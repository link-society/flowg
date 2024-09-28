package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
)

type SaveUserRequest struct {
	User     string   `path:"user" minLength:"1"`
	Roles    []string `json:"roles"`
	Password string   `json:"password"`
}

type SaveUserResponse struct {
	Success bool `json:"success"`
}

func SaveUserUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_ACLS,
			func(
				ctx context.Context,
				req SaveUserRequest,
				resp *SaveUserResponse,
			) error {
				user := models.User{
					Name:  req.User,
					Roles: req.Roles,
				}

				err := authStorage.SaveUser(ctx, user, req.Password)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save user",
						slog.String("channel", "api"),
						slog.String("user", req.User),
						slog.String("error", err.Error()),
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
