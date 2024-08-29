package admin

import (
	"net/http"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/admin/controllers"
)

func Application(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	roleSys := auth.NewRoleSystem(authDb)
	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /web/admin/{$}",
		controllers.Index(roleSys, userSys),
	)

	mux.HandleFunc(
		"GET /web/admin/roles/new/{$}",
		controllers.DisplayRoleCreateForm(userSys),
	)
	mux.HandleFunc(
		"POST /web/admin/roles/new/{$}",
		controllers.ProcessRoleCreateForm(roleSys, userSys),
	)
	mux.HandleFunc(
		"POST /web/admin/roles/delete/{name}/{$}",
		controllers.ProcessRoleDeleteAction(roleSys, userSys),
	)

	mux.HandleFunc(
		"GET /web/admin/users/new/{$}",
		controllers.DisplayUserCreateForm(roleSys, userSys),
	)
	mux.HandleFunc(
		"POST /web/admin/users/new/{$}",
		controllers.ProcessUserCreateForm(roleSys, userSys),
	)
	mux.HandleFunc(
		"POST /web/admin/users/delete/{name}/{$}",
		controllers.ProcessUserDeleteAction(userSys),
	)

	return mux
}
