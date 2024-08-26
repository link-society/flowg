package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/logstorage"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/templates/views"
)

func MainController(
	authDb *auth.Database,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/{$}", func(w http.ResponseWriter, r *http.Request) {
		streamCount := 0
		transformerCount := 0
		pipelineCount := 0

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

			notifications = append(notifications, "❌ Could not fetch user permissions")
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
			permissions,
			notifications,
		))

		h.ServeHTTP(w, r)
	})

	return mux
}
