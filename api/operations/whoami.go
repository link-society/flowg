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
	"link-society.com/flowg/internal/models"

	authStorage "link-society.com/flowg/internal/storage/auth"
)

// WhoamiDeps lists the dependencies of [NewWhoamiUsecase].
type WhoamiDeps struct {
	fx.In

	AuthStorage authStorage.Storage
}

// WhoamiRequest is empty: the caller is identified by their credentials.
type WhoamiRequest struct{}

// WhoamiResponse describes the currently authenticated user.
type WhoamiResponse struct {
	// Success reports whether the profile was returned.
	Success bool `json:"success"`
	// User is the authenticated account and its assigned roles.
	User *models.User `json:"user"`
	// Permissions is the effective permission set derived from the user's roles.
	Permissions models.Permissions `json:"permissions"`
}

// NewWhoamiUsecase returns the profile and effective permissions of the caller.
//
// It lets a client discover who it is authenticated as and what it is allowed
// to do, typically to drive UI affordances. It requires authentication but no
// particular permission.
func NewWhoamiUsecase(deps WhoamiDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req WhoamiRequest,
			resp *WhoamiResponse,
		) error {
			resp.User = auth.GetContextUser(ctx)

			scopes, err := deps.AuthStorage.ListUserScopes(ctx, resp.User.Name)
			if err != nil {
				logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewWhoamiUsecase,
		http.MethodGet,
		"/api/v1/auth/whoami",
	)
}
