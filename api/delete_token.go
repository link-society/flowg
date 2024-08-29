package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type DeleteTokenRequest struct {
	TokenUUID string `path:"token-uuid" format:"uuid"`
}

type DeleteTokenResponse struct {
	Success bool `json:"success"`
}

func DeleteTokenUsecase(authDb *auth.Database) usecase.Interactor {
	tokenSys := auth.NewTokenSystem(authDb)

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req DeleteTokenRequest,
			resp *DeleteTokenResponse,
		) error {
			user := auth.GetContextUser(ctx)

			err := tokenSys.DeleteToken(user.Name, req.TokenUUID)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to delete token",
					"channel", "api",
					"user", user.Name,
					"token-uuid", req.TokenUUID,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true

			return nil
		},
	)

	u.SetName("delete_token")
	u.SetTitle("Delete Token")
	u.SetDescription("Delete Personal Access Token UUIDs for the current user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.Internal)

	return u
}
