package controllers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func TransformersController(pipelinesManager *pipelines.Manager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/transformers/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/transformers/new", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET /web/transformers/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		h := templ.Handler(views.Transformers(
			views.TransformersProps{
				Transformers:       transformers,
				CurrentTransformer: "",
				Code:               ".",
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	return mux
}
