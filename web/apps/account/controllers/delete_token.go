package controllers

import (
	"log/slog"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils/htmx"
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

			htmx.Reswap(w, "none")

			notifications = append(notifications, "&#10060; Could not delete token")
		} else {
			htmx.Reswap(w, "delete")
			htmx.Retarget(w, "tr[data-token='"+tokenUUID+"']")

			notifications = append(notifications, "&#9989; Token deleted")
		}

		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: notifications,
			},
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
