package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func PipelinesController(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/pipelines/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/pipelines/new", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("GET /web/pipelines/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
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
			notifications = append(notifications, "&#10060; Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: "",
				Flow:            "null",
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("POST /web/pipelines/new/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		pipelineName := ""
		pipelineFlow := "."

		err = r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not parse form")
		} else {
			pipelineName = r.FormValue("name")
			pipelineFlow = r.FormValue("flow")

			if permissions.CanEditPipelines {
				if pipelineName == "" {
					notifications = append(notifications, "&#10060; Pipeline name is required")
				}

				if pipelineFlow == "" {
					notifications = append(notifications, "&#10060; Pipeline flow is required")
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

						notifications = append(notifications, "&#10060; Could not save pipeline")
					} else {
						notifications = append(notifications, "&#9989; Pipeline saved")
					}
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
			notifications = append(notifications, "&#10060; Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			permissions,
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

		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
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
			notifications = append(notifications, "&#10060; Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			permissions,
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

		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		err = r.ParseForm()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error parsing form",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not parse form")
		} else {
			pipelineName = r.FormValue("name")
			pipelineFlow = r.FormValue("flow")

			if !permissions.CanEditPipelines {
				if pipelineName == "" {
					notifications = append(notifications, "&#10060; Pipeline name is required")
				}

				if pipelineFlow == "" {
					notifications = append(notifications, "&#10060; Pipeline flow is required")
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

						notifications = append(notifications, "&#10060; Could not save pipeline")
					} else {
						notifications = append(notifications, "&#9989; Pipeline saved")
					}
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
			notifications = append(notifications, "&#10060; Could not fetch pipelines")
		}

		h := templ.Handler(views.Pipelines(
			views.PipelinesProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/pipelines/delete/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		pipelineName := r.PathValue("name")

		permissions := auth.Permissions{}
		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewPipelines {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		} else if !permissions.CanEditPipelines {
			http.Redirect(w, r, fmt.Sprintf("/web/pipelines/edit/%s", pipelineName), http.StatusSeeOther)
			return
		}

		err = pipelinesManager.DeletePipelineFlow(pipelineName)
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
