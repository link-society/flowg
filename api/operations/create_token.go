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

	"link-society.com/flowg/internal/storage"
)

// CreateTokenDeps lists the dependencies of [NewCreateTokenUsecase].
type CreateTokenDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// CreateTokenRequest is empty: the token is issued for the calling user.
type CreateTokenRequest struct{}

// CreateTokenResponse carries the newly issued personal access token.
type CreateTokenResponse struct {
	// Success reports whether the token was created.
	Success bool `json:"success"`
	// Token is the secret value, returned only once at creation time.
	Token string `json:"token"`
	// TokenUUID identifies the token for later listing or deletion.
	TokenUUID string `json:"token_uuid"`
}

// NewCreateTokenUsecase issues a new personal access token for the calling user.
//
// The secret value is returned only in this response and cannot be retrieved
// later; only its UUID is persisted for management. It requires authentication
// but no particular permission.
func NewCreateTokenUsecase(deps CreateTokenDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req CreateTokenRequest,
			resp *CreateTokenResponse,
		) error {
			user := auth.GetContextUser(ctx)

			token, tokenUuid, err := deps.AuthStorage.CreateToken(ctx, user.Name)
			if err != nil {
				logger.ErrorContext(
					ctx,
					"Failed to create token",
					slog.String("user", user.Name),
					slog.String("error", err.Error()),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true
			resp.Token = token
			resp.TokenUUID = tokenUuid

			return nil
		},
	)

	u.SetName("create_token")
	u.SetTitle("Create Token")
	u.SetDescription("Create a new Personal Access Token for the current user")
	u.SetTags("acls")

	u.SetExpectedErrors(status.NotFound, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewCreateTokenUsecase,
		http.MethodPost,
		"/api/v1/token",
	)
}
