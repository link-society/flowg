package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/storage/auth"
)

type ListTokensRequest struct{}

type ListTokensResponse struct {
	Success    bool     `json:"success"`
	TokenUUIDs []string `json:"token_uuids"`
}

func ListTokensUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ListTokensRequest,
			resp *ListTokensResponse,
		) error {
			user := apiUtils.GetContextUser(ctx)

			tokenUUIDs, err := authStorage.ListTokens(ctx, user.Name)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to list tokens",
					slog.String("channel", "api"),
					slog.String("user", user.Name),
					slog.String("error", err.Error()),
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
