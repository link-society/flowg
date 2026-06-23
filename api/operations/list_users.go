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

// ListUsersDeps lists the dependencies of [NewListUsersUsecase].
type ListUsersDeps struct {
	fx.In

	AuthStorage authStorage.Storage
}

// ListUsersRequest is empty: listing users takes no parameters.
type ListUsersRequest struct{}

// ListUsersResponse carries every known user with its roles.
type ListUsersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Users holds every account and its assigned roles.
	Users []models.User `json:"users"`
}

// NewListUsersUsecase enumerates all user accounts with their roles.
//
// Callers must have the read-ACLs permission.
func NewListUsersUsecase(deps ListUsersDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListUsersRequest,
				resp *ListUsersResponse,
			) error {
				users, err := deps.AuthStorage.ListUsers(ctx)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to list users",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Users = users
				return nil
			},
		),
	)

	u.SetName("list_users")
	u.SetTitle("List Users")
	u.SetDescription("List known users")
	u.SetTags("acls")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListUsersUsecase,
		http.MethodGet,
		"/api/v1/users",
	)
}
