package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/account/templates/components"
)

func CreateToken(
	tokenSys *auth.TokenSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		success := false
		notifications := []string{}
		user := auth.GetContextUser(r.Context())

		token, err := tokenSys.CreateToken(user.Name)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error generating token",
				"channel", "web",
				"user", user.Name,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not create token")
		} else {
			success = true
		}

		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: notifications,
			},
		}

		if success {
			trigger.ModalOpenEvent = &htmx.ModalOpenEvent{}
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
