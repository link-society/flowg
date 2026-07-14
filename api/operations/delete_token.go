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
	"link-society.com/flowg/api/schemas"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// DeleteTokenDeps lists the dependencies of [NewDeleteTokenUsecase].
type DeleteTokenDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// NewDeleteTokenUsecase revokes one of the calling user's personal access tokens.
//
// A user may only delete their own tokens. It requires authentication but no
// particular permission.
func NewDeleteTokenUsecase(deps DeleteTokenDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req schemas.DeleteTokenRequest,
			resp *schemas.DeleteTokenResponse,
		) error {
			user := auth.GetContextUser(ctx)

			err := deps.AuthStorage.DeleteToken(ctx, user.Name, req.TokenUUID)
			if err != nil {
				logger.ErrorContext(
					ctx,
					"Failed to delete token",
					slog.String("user", user.Name),
					slog.String("token-uuid", req.TokenUUID),
					slog.String("error", err.Error()),
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

func init() {
	routing.RegisterOperation(
		NewDeleteTokenUsecase,
		http.MethodDelete,
		"/api/v1/tokens/{token-uuid}",
	)
}
