package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/apps/dashboard/templates/views"
)

func Page(
	userSys *auth.UserSystem,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		streamCount := 0
		transformerCount := 0
		pipelineCount := 0

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

		streamList, err := logDb.ListStreams()
		if err == nil {
			streamCount = len(streamList)
		} else {
			slog.ErrorContext(
				r.Context(),
				"error listing streams",
				"channel", "web",
				"error", err.Error(),
			)
			notifications = append(notifications, "&#10060; Could not fetch streams")
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
			notifications = append(notifications, "&#10060; Could not fetch transformers")
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
			notifications = append(notifications, "&#10060; Could not fetch pipelines")
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				StreamCount:      streamCount,
				TransformerCount: transformerCount,
				PipelineCount:    pipelineCount,

				Permissions:   permissions,
				Notifications: notifications,
			},
		))

		h.ServeHTTP(w, r)
	}
}
