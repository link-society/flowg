package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type WhoamiRequest struct{}
type WhoamiResponse struct {
	Success     bool             `json:"success"`
	User        *auth.User       `json:"user"`
	Permissions auth.Permissions `json:"permissions"`
}

func WhoamiUsecase(authDb *auth.Database) usecase.Interactor {
	userSys := auth.NewUserSystem(authDb)

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req WhoamiRequest,
			resp *WhoamiResponse,
		) error {
			resp.User = auth.GetContextUser(ctx)

			scopes, err := userSys.ListUserScopes(resp.User.Name)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to fetch user scopes",
					"channel", "api",
					"error", err.Error(),
				)
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Permissions = auth.PermissionsFromScopes(scopes)
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
