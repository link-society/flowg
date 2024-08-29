package onboarding

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"

	"link-society.com/flowg/web/apps/onboarding/controllers"
)

func Application(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /auth/login/{$}",
		controllers.LoginForm(userSys),
	)
	mux.HandleFunc(
		"POST /auth/login/{$}",
		controllers.LoginAction(userSys),
	)
	mux.HandleFunc(
		"GET /auth/logout/{$}",
		controllers.LogoutAction(),
	)

	return mux
}
