package controllers

import (
	"log/slog"

	"encoding/json"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"

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

		trigger := map[string]interface{}{
			"htmx-custom-toast": map[string]interface{}{
				"messages": notifications,
			},
		}

		if success {
			trigger["htmx-custom-modal-open"] = map[string]interface{}{}
		}

		triggerData, err := json.Marshal(trigger)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error marshalling trigger",
				"channel", "web",
				"error", err.Error(),
			)
		} else {
			w.Header().Add("HX-Trigger", string(triggerData))
		}

		h := templ.Handler(components.TokenViewer(
			components.TokenViewerProps{
				Token: token,
			},
		))
		h.ServeHTTP(w, r)
	}
}
