package controllers

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessStreamDeleteAction(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		trigger := htmx.Trigger{}
		streamName := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanEditStreams {
			webutils.NotifyError(r.Context(), "You do not have permission to delete streams")
			goto response
		}

		if err := metaSys.DeleteStream(r.Context(), streamName); err != nil {
			webutils.LogError(r.Context(), "Failed to delete stream", err)
			webutils.NotifyError(r.Context(), "Failed to delete stream")
			goto response
		}

		w.Header().Add("HX-Redirect", "/web/storage/")

	response:
		trigger.ToastEvent = &htmx.ToastEvent{
			Messages: webutils.Notifications(r.Context()),
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
