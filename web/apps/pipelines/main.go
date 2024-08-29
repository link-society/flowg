package pipelines

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"

	"link-society.com/flowg/web/apps/pipelines/controllers"
)

func Application(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /web/pipelines/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/web/pipelines/new", http.StatusPermanentRedirect)
		},
	)

	mux.HandleFunc(
		"GET /web/pipelines/new/{$}",
		controllers.PageNew(userSys, pipelinesManager),
	)
	mux.HandleFunc(
		"POST /web/pipelines/new/{$}",
		controllers.ProcessNewSaveAction(userSys, pipelinesManager),
	)

	mux.HandleFunc(
		"GET /web/pipelines/edit/{name}/{$}",
		controllers.PageEdit(userSys, pipelinesManager),
	)
	mux.HandleFunc(
		"POST /web/pipelines/edit/{name}/{$}",
		controllers.ProcessEditSaveAction(userSys, pipelinesManager),
	)

	mux.HandleFunc(
		"GET /web/pipelines/delete/{name}/{$}",
		controllers.ProcessDeleteAction(userSys, pipelinesManager),
	)

	return mux
}
