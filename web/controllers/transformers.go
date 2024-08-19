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

	mux.HandleFunc("POST /web/transformers/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}
		transformerName := ""
		transformerCode := "."

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
			transformerName = r.FormValue("name")
			transformerCode = r.FormValue("code")

			if transformerName == "" {
				notifications = append(notifications, "❌ Transformer name is required")
			}

			if transformerCode == "" {
				notifications = append(notifications, "❌ Transformer code is required")
			}

			if transformerName != "" && transformerCode != "" {
				err = pipelinesManager.SaveTransformerScript(transformerName, transformerCode)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving transformer script",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "❌ Could not save transformer script")
				} else {
					notifications = append(notifications, "✅ Transformer script saved")
				}
			}
		}

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
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/transformers/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")
		transformerCode, err := pipelinesManager.GetTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
			return
		}

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
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/transformers/edit/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")
		transformerCode, err := pipelinesManager.GetTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
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
			transformerName = r.FormValue("name")
			transformerCode = r.FormValue("code")

			if transformerName == "" {
				notifications = append(notifications, "❌ Transformer name is required")
			}

			if transformerCode == "" {
				notifications = append(notifications, "❌ Transformer code is required")
			}

			if transformerName != "" && transformerCode != "" {
				err = pipelinesManager.SaveTransformerScript(transformerName, transformerCode)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving transformer script",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "❌ Could not save transformer script")
				} else {
					notifications = append(notifications, "✅ Transformer script saved")
				}
			}
		}

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
				CurrentTransformer: transformerName,
				Code:               transformerCode,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	return mux
}
