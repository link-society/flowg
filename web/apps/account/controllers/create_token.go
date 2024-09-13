package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/account/templates/components"
)

func CreateToken(
	tokenSys *auth.TokenSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trigger htmx.Trigger

		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))

		user := auth.GetContextUser(r.Context())
		token, _, err := tokenSys.CreateToken(user.Name)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to create token", err)
			webutils.NotifyError(r.Context(), "Could not create token")
			goto response
		}

		trigger.ModalOpenEvent = &htmx.ModalOpenEvent{}

	response:
		trigger.ToastEvent = &htmx.ToastEvent{
			Messages: webutils.Notifications(r.Context()),
		}

		trigger.Write(r.Context(), w)

		h := templ.Handler(components.TokenViewer(
			components.TokenViewerProps{
				Token: token,
			},
		))
		h.ServeHTTP(w, r)
	}
}
