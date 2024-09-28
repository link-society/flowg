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

type WhoamiRequest struct{}
type WhoamiResponse struct {
	Success     bool               `json:"success"`
	User        *models.User       `json:"user"`
	Permissions models.Permissions `json:"permissions"`
}

func WhoamiUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req WhoamiRequest,
			resp *WhoamiResponse,
		) error {
			resp.User = apiUtils.GetContextUser(ctx)

			scopes, err := authStorage.ListUserScopes(ctx, resp.User.Name)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to fetch user scopes",
					slog.String("channel", "api"),
					slog.String("error", err.Error()),
				)
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Permissions = models.PermissionsFromScopes(scopes)
			return nil
		},
	)

	u.SetName("whoami")
	u.SetTitle("Fetch current profile")
	u.SetDescription("Fetch the profile of the currently authenticated user")
	u.SetTags("auth")

	u.SetExpectedErrors(status.Internal)

	return u
}
