package controllers

import (
	"fmt"
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessStreamCreateForm(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		trigger := htmx.Trigger{}
		var streamName string

		if !webutils.Permissions(r.Context()).CanEditStreams {
			webutils.NotifyError(r.Context(), "You do not have permission to edit streams")
			goto response
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Failed to parse form")
			goto response
		}

		streamName = r.FormValue("stream_name")
		if err := metaSys.ConfigureStream(streamName, logstorage.StreamConfig{}); err != nil {
			webutils.LogError(r.Context(), "Failed to create stream", err)
			webutils.NotifyError(r.Context(), "Failed to create stream")
			goto response
		}

		w.Header().Add("HX-Redirect", fmt.Sprintf("/web/storage/edit/%s/", streamName))

	response:
		trigger.ToastEvent = &htmx.ToastEvent{
			Messages: webutils.Notifications(r.Context()),
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
