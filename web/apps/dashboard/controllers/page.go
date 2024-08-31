package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/data/pipelines"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/dashboard/templates/views"
)

func Page(
	userSys *auth.UserSystem,
	logDb *logstorage.Storage,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	metaSys := logstorage.NewMetaSystem(logDb)

	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		streamCount := 0
		transformerCount := 0
		pipelineCount := 0

		streamList, err := metaSys.ListStreams()
		if err == nil {
			streamCount = len(streamList)
		} else {
			webutils.LogError(r.Context(), "Failed to fetch streams", err)
			webutils.NotifyError(r.Context(), "Could not fetch streams")
		}

		transformerList, err := pipelinesManager.ListTransformers()
		if err == nil {
			transformerCount = len(transformerList)
		} else {
			webutils.LogError(r.Context(), "Failed to fetch transformers", err)
			webutils.NotifyError(r.Context(), "Could not fetch transformers")
		}

		pipelineList, err := pipelinesManager.ListPipelines()
		if err == nil {
			pipelineCount = len(pipelineList)
		} else {
			webutils.LogError(r.Context(), "Failed to fetch pipelines", err)
			webutils.NotifyError(r.Context(), "Could not fetch pipelines")
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				StreamCount:      streamCount,
				TransformerCount: transformerCount,
				PipelineCount:    pipelineCount,

				Permissions:   webutils.Permissions(r.Context()),
				Notifications: webutils.Notifications(r.Context()),
			},
		))

		h.ServeHTTP(w, r)
	}
}
