package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
)

type ListTokensRequest struct{}

type ListTokensResponse struct {
	Success    bool     `json:"success"`
	TokenUUIDs []string `json:"token-uuids"`
}

func ListTokensUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ListTokensRequest,
			resp *ListTokensResponse,
		) error {
			username := auth.GetContextUser(ctx)

			tokenUUIDs, err := authDb.ListPersonalAccessTokens(username)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to list tokens",
					"channel", "api",
					"user", username,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.TokenUUIDs = tokenUUIDs

			return nil
		},
	)

	u.SetName("list_tokens")
	u.SetTitle("List Tokens")
	u.SetDescription("List Personal Access Token UUIDs for the current user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.Internal)

	return u
}
