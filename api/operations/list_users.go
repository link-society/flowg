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

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ListUsersDeps lists the dependencies of [NewListUsersUsecase].
type ListUsersDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
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
				req schemas.ListUsersRequest,
				resp *schemas.ListUsersResponse,
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
