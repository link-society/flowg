package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/pipelines"
	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/web/templates/views"
)

func MainController(
	db *storage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/{$}", func(w http.ResponseWriter, r *http.Request) {
		streamCount := 0
		transformerCount := 0
		pipelineCount := 0

		streamList, err := db.ListStreams()
		if err == nil {
			streamCount = len(streamList)
		}

		transformerList, err := pipelinesManager.ListTransformers()
		if err == nil {
			transformerCount = len(transformerList)
		}

		pipelineList, err := pipelinesManager.ListPipelines()
		if err == nil {
			pipelineCount = len(pipelineList)
		}

		h := templ.Handler(views.Dashboard(
			streamCount,
			transformerCount,
			pipelineCount,
		))

		h.ServeHTTP(w, r)
	})

	return mux
}
