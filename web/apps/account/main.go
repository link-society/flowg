package account

import (
	"net/http"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/account/controllers"
)

func Application(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	tokenSys := auth.NewTokenSystem(authDb)

	mux.HandleFunc(
		"GET /web/account/{$}",
		controllers.Index(userSys, tokenSys),
	)

	mux.HandleFunc(
		"POST /web/account/change-password/{$}",
		controllers.ChangePassword(userSys),
	)

	mux.HandleFunc(
		"POST /web/account/token/new/{$}",
		controllers.CreateToken(tokenSys),
	)
	mux.HandleFunc(
		"POST /web/account/token/delete/{tokenUUID}/{$}",
		controllers.DeleteToken(tokenSys),
	)

	return mux
}
