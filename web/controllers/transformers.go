package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func TransformersController(pipelinesManager *pipelines.Manager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/transformers/{$}", func(w http.ResponseWriter, r *http.Request) {
		h := templ.Handler(views.Transformers())
		h.ServeHTTP(w, r)
	})

	return mux
}
