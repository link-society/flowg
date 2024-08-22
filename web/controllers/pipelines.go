package controllers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func PipelinesController(pipelinesManager *pipelines.Manager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/pipelines/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/pipelines/new", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET /web/pipelines/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		pipelines, err := pipelinesManager.ListPipelines()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing pipelines",
				"channel", "web",
				"error", err.Error(),
			)

			pipelines = []string{}
			notifications = append(notifications, "❌ Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: "",
				Flow:            "null",
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/pipelines/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}
		pipelineName := ""
		pipelineFlow := "."

		err := r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not parse form")
		} else {
			pipelineName = r.FormValue("name")
			pipelineFlow = r.FormValue("flow")

			if pipelineName == "" {
				notifications = append(notifications, "❌ Pipeline name is required")
			}

			if pipelineFlow == "" {
				notifications = append(notifications, "❌ Pipeline flow is required")
			}

			if pipelineName != "" && pipelineFlow != "" {
				err = pipelinesManager.SavePipelineFlow(pipelineName, pipelineFlow)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving pipeline flow",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "❌ Could not save pipeline")
				} else {
					notifications = append(notifications, "✅ Pipeline saved")
				}
			}
		}

		pipelines, err := pipelinesManager.ListPipelines()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing pipelines",
				"channel", "web",
				"error", err.Error(),
			)

			pipelines = []string{}
			notifications = append(notifications, "❌ Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/pipelines/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("name")
		pipelineFlow, err := pipelinesManager.GetPipelineFlow(pipelineName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting pipeline flow",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
			return
		}

		notifications := []string{}

		pipelines, err := pipelinesManager.ListPipelines()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing pipelines",
				"channel", "web",
				"error", err.Error(),
			)

			pipelines = []string{}
			notifications = append(notifications, "❌ Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/pipelines/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("name")
		pipelineFlow, err := pipelinesManager.GetPipelineFlow(pipelineName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting pipeline flow",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
			return
		}

		notifications := []string{}

		err = r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not parse form")
		} else {
			pipelineName = r.FormValue("name")
			pipelineFlow = r.FormValue("flow")

			if pipelineName == "" {
				notifications = append(notifications, "❌ Pipeline name is required")
			}

			if pipelineFlow == "" {
				notifications = append(notifications, "❌ Pipeline flow is required")
			}

			if pipelineName != "" && pipelineFlow != "" {
				err = pipelinesManager.SavePipelineFlow(pipelineName, pipelineFlow)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving pipeline flow",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "❌ Could not save pipeline")
				} else {
					notifications = append(notifications, "✅ Pipeline saved")
				}
			}
		}

		pipelines, err := pipelinesManager.ListPipelines()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing pipeline",
				"channel", "web",
				"error", err.Error(),
			)

			pipelines = []string{}
			notifications = append(notifications, "❌ Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/pipelines/delete/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("name")
		err := pipelinesManager.DeletePipelineFlow(pipelineName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error deleting pipeline",
				"channel", "web",
				"error", err.Error(),
			)
		}

		http.Redirect(w, r, "/web/pipelines/new", http.StatusTemporaryRedirect)
	})

	return mux
}
