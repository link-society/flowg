package pipelines

import (
	"log/slog"
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"

	"link-society.com/flowg/web/apps/pipelines/controllers"
)

func Application(
	authDb *auth.Database,
	configStorage *config.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	pipelineSys := config.NewPipelineSystem(configStorage)

	mux.HandleFunc(
		"GET /web/pipelines/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			switch pipelines, err := pipelineSys.List(); {
			case err == nil && len(pipelines) > 0:
				url := "/web/pipelines/edit/" + pipelines[0] + "/"
				http.Redirect(w, r, url, http.StatusSeeOther)

			case err != nil:
				slog.ErrorContext(
					r.Context(),
					"Failed to list pipelines",
					"channel", "web",
					"error", err,
				)
				fallthrough

			default:
				http.Redirect(w, r, "/web/pipelines/new", http.StatusSeeOther)
			}
		},
	)

	mux.HandleFunc(
		"GET /web/pipelines/new/{$}",
		controllers.PageNew(userSys, pipelineSys),
	)
	mux.HandleFunc(
		"POST /web/pipelines/new/{$}",
		controllers.ProcessNewSaveAction(userSys, pipelineSys),
	)

	mux.HandleFunc(
		"GET /web/pipelines/edit/{name}/{$}",
		controllers.PageEdit(userSys, pipelineSys),
	)
	mux.HandleFunc(
		"POST /web/pipelines/edit/{name}/{$}",
		controllers.ProcessEditSaveAction(userSys, pipelineSys),
	)

	mux.HandleFunc(
		"GET /web/pipelines/delete/{name}/{$}",
		controllers.ProcessDeleteAction(userSys, pipelineSys),
	)

	return mux
}
