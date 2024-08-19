package controllers

import (
	"log/slog"

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

		notifications := []string{}

		streamList, err := db.ListStreams()
		if err == nil {
			streamCount = len(streamList)
		} else {
			slog.ErrorContext(
				r.Context(),
				"error listing streams",
				"channel", "web",
				"error", err.Error(),
			)
			notifications = append(notifications, "❌ Could not fetch streams")
		}

		transformerList, err := pipelinesManager.ListTransformers()
		if err == nil {
			transformerCount = len(transformerList)
		} else {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)
			notifications = append(notifications, "❌ Could not fetch transformers")
		}

		pipelineList, err := pipelinesManager.ListPipelines()
		if err == nil {
			pipelineCount = len(pipelineList)
		} else {
			slog.ErrorContext(
				r.Context(),
				"error listing pipelines",
				"channel", "web",
				"error", err.Error(),
			)
			notifications = append(notifications, "❌ Could not fetch pipelines")
		}

		h := templ.Handler(views.Dashboard(
			views.DashboardProps{
				StreamCount:      streamCount,
				TransformerCount: transformerCount,
				PipelineCount:    pipelineCount,
			},
			notifications,
		))

		h.ServeHTTP(w, r)
	})

	return mux
}
