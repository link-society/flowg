package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"

	"link-society.com/flowg/web/apps/pipelines/templates/views"
)

func ProcessNewSaveAction(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := userSys.ListUserScopes(user.Name)
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

		h := templ.Handler(views.Page(
			views.PageProps{
				Pipelines:       pipelines,
				CurrentPipeline: pipelineName,
				Flow:            pipelineFlow,

				Permissions:   permissions,
				Notifications: notifications,
			},
		))
		h.ServeHTTP(w, r)
	}
}
