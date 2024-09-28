package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
)

type ListUsersRequest struct{}
type ListUsersResponse struct {
	Success bool          `json:"success"`
	Users   []models.User `json:"users"`
}

func ListUsersUsecase(authStorage *auth.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_ACLS,
			func(
				ctx context.Context,
				req ListUsersRequest,
				resp *ListUsersResponse,
			) error {
				users, err := authStorage.ListUsers(ctx)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list users",
						slog.String("channel", "api"),
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
