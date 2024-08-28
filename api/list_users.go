package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/auth"
)

type ListUsersRequest struct{}
type ListUsersResponse struct {
	Success bool         `json:"success"`
	Users   []*auth.User `json:"Users"`
}

func ListUsersUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiMiddleware(
			authDb,
			auth.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListUsersRequest,
				resp *ListUsersResponse,
			) error {
				usernames, err := authDb.ListUsers()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list users",
						"channel", "api",
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				users := make([]*auth.User, 0, len(usernames))

				for i, username := range usernames {
					user, err := authDb.GetUser(username)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to get user",
							"channel", "api",
							"user", username,
							"error", err.Error(),
						)

						resp.Success = false
						return status.Wrap(err, status.Internal)
					}

					users[i] = user
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
