package controllers

import (
	"log/slog"

	"encoding/json"

	"net/http"

	"link-society.com/flowg/internal/auth"
)

func DeleteToken(
	tokenSys *auth.TokenSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.GetContextUser(r.Context())
		notifications := []string{}

		tokenUUID := r.PathValue("tokenUUID")
		err := tokenSys.DeleteToken(user.Name, tokenUUID)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error deleting token",
				"channel", "web",
				"user", user.Name,
				"tokenUUID", tokenUUID,
				"error", err.Error(),
			)

			w.Header().Add("HX-Reswap", "none")

			notifications = append(notifications, "&#10060; Could not delete token")
		} else {
			w.Header().Add("HX-Reswap", "delete")
			w.Header().Add("HX-Retarget", "tr[data-token='"+tokenUUID+"']")

			notifications = append(notifications, "&#9989; Token deleted")
		}

		trigger := map[string]interface{}{
			"htmx-custom-toast": map[string]interface{}{
				"messages": notifications,
			},
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

		w.WriteHeader(http.StatusOK)
	}
}
