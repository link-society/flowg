package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/storage/templates/components"
)

func ProcessIndexAddAction(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanEditStreams {
			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to edit streams"))
			return
		}

		streamName := r.PathValue("name")

		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}
		trigger.Write(r.Context(), w)

		h := templ.Handler(components.NewIndexFieldForm(
			components.NewIndexFieldFormProps{
				StreamName: streamName,
			},
		))
		h.ServeHTTP(w, r)
	}
}
