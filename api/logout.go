package api

import (
	"context"

	"github.com/swaggest/usecase"

	"link-society.com/flowg/internal/data/auth"
)

type LogoutRequest struct{}

type LogoutResponse struct {
	Success   bool   `json:"success"`
	SessionID string `cookie:"session_id,httponly,max-age=-1,path=/"`
}

func LogoutUsecase(authDb *auth.Database) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req LogoutRequest,
			resp *LogoutResponse,
		) error {
			resp.SessionID = ""
			return nil
		},
	)

	u.SetName("logout")
	u.SetTitle("Sign out")
	u.SetDescription("Delete Session cookie")
	u.SetTags("auth")

	return u
}
