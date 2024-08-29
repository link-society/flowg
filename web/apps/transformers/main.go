package transformers

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"

	"link-society.com/flowg/web/apps/transformers/controllers"
)

func Application(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)

	mux.HandleFunc(
		"GET /web/transformers/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/web/transformers/new", http.StatusPermanentRedirect)
		},
	)

	mux.HandleFunc(
		"GET /web/transformers/new/{$}",
		controllers.PageNew(userSys, pipelinesManager),
	)
	mux.HandleFunc(
		"POST /web/transformers/new/{$}",
		controllers.ProcessNewSaveAction(userSys, pipelinesManager),
	)

	mux.HandleFunc(
		"GET /web/transformers/edit/{name}/{$}",
		controllers.PageEdit(userSys, pipelinesManager),
	)
	mux.HandleFunc(
		"POST /web/transformers/edit/{name}/{$}",
		controllers.ProcessEditSaveAction(userSys, pipelinesManager),
	)

	mux.HandleFunc(
		"GET /web/transformers/delete/{name}/{$}",
		controllers.ProcessDeleteAction(userSys, pipelinesManager),
	)

	return mux
}
