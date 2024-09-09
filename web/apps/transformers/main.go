package transformers

import (
	"log/slog"
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"

	"link-society.com/flowg/web/apps/transformers/controllers"
)

func Application(
	authDb *auth.Database,
	configStorage *config.Storage,
) http.Handler {
	mux := http.NewServeMux()

	userSys := auth.NewUserSystem(authDb)
	transformerSys := config.NewTransformerSystem(configStorage)

	mux.HandleFunc(
		"GET /web/transformers/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			switch transformers, err := transformerSys.List(); {
			case err == nil && len(transformers) > 0:
				url := "/web/transformers/edit/" + transformers[0] + "/"
				http.Redirect(w, r, url, http.StatusSeeOther)

			case err != nil:
				slog.ErrorContext(
					r.Context(),
					"Failed to list transformers",
					"channel", "web",
					"error", err,
				)
				fallthrough

			default:
				http.Redirect(w, r, "/web/transformers/new", http.StatusSeeOther)
			}
		},
	)

	mux.HandleFunc(
		"GET /web/transformers/new/{$}",
		controllers.PageNew(userSys, transformerSys),
	)
	mux.HandleFunc(
		"POST /web/transformers/new/{$}",
		controllers.ProcessNewSaveAction(userSys, transformerSys),
	)

	mux.HandleFunc(
		"GET /web/transformers/edit/{name}/{$}",
		controllers.PageEdit(userSys, transformerSys),
	)
	mux.HandleFunc(
		"POST /web/transformers/edit/{name}/{$}",
		controllers.ProcessEditSaveAction(userSys, transformerSys),
	)

	mux.HandleFunc(
		"GET /web/transformers/delete/{name}/{$}",
		controllers.ProcessDeleteAction(userSys, transformerSys),
	)

	return mux
}
