package controllers

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func DeleteToken(
	tokenSys *auth.TokenSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))

		user := auth.GetContextUser(r.Context())
		tokenUUID := r.PathValue("tokenUUID")
		err := tokenSys.DeleteToken(user.Name, tokenUUID)
		if err != nil {
			webutils.LogError(
				r.Context(),
				"Failed to delete token", err,
				"tokenUUID", tokenUUID,
			)
			htmx.Reswap(w, "none")

			webutils.NotifyError(r.Context(), "Could not delete token")
			goto response
		}

		htmx.Reswap(w, "delete")
		htmx.Retarget(w, "tr[data-token='"+tokenUUID+"']")
		webutils.NotifyInfo(r.Context(), "Token deleted")

	response:
		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
