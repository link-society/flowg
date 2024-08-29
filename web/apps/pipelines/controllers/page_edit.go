package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/apps/pipelines/templates/views"
)

func PageEdit(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
