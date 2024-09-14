package api

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
)

type WhoamiRequest struct{}
type WhoamiResponse struct {
	Success bool       `json:"success"`
	User    *auth.User `json:"user"`
}

func WhoamiUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req WhoamiRequest,
			resp *WhoamiResponse,
		) error {
			resp.Success = true
			resp.User = auth.GetContextUser(ctx)
			return nil
		},
	)

	u.SetName("whoami")
	u.SetTitle("Fetch current profile")
	u.SetDescription("Fetch the profile of the currently authenticated user")
	u.SetTags("auth")

	u.SetExpectedErrors(status.PermissionDenied)

	return u
}
