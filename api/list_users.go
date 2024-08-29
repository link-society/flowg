package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/data/auth"
)

type ListUsersRequest struct{}
type ListUsersResponse struct {
	Success bool        `json:"success"`
	Users   []auth.User `json:"Users"`
}

func ListUsersUsecase(authDb *auth.Database) usecase.Interactor {
	userSys := auth.NewUserSystem(authDb)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListUsersRequest,
				resp *ListUsersResponse,
			) error {
				users, err := userSys.ListUsers()
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
