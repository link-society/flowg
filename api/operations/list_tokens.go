package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"

	authStorage "link-society.com/flowg/internal/storage/auth"
)

// ListTokensDeps lists the dependencies of [NewListTokensUsecase].
type ListTokensDeps struct {
	fx.In

	AuthStorage authStorage.Storage
}

// ListTokensRequest is empty: tokens are listed for the calling user.
type ListTokensRequest struct{}

// ListTokensResponse carries the identifiers of the caller's tokens.
type ListTokensResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// TokenUUIDs identifies each of the caller's tokens; the secret values are
	// never returned.
	TokenUUIDs []string `json:"token_uuids"`
}

// NewListTokensUsecase enumerates the identifiers of the calling user's personal
// access tokens.
//
// Only the UUIDs are returned; the secret values are never recoverable after
// creation. It requires authentication but no particular permission.
func NewListTokensUsecase(deps ListTokensDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ListTokensRequest,
			resp *ListTokensResponse,
		) error {
			user := auth.GetContextUser(ctx)

			tokenUUIDs, err := deps.AuthStorage.ListTokens(ctx, user.Name)
			if err != nil {
				logger.ErrorContext(
					ctx,
					"Failed to list tokens",
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

func init() {
	routing.RegisterOperation(
		NewListTokensUsecase,
		http.MethodGet,
		"/api/v1/tokens",
	)
}
