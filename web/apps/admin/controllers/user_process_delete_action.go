package controllers

import (
	"strconv"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessUserDeleteAction(
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		username := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanEditACLs {
			htmx.Reswap(w, "none")
			webutils.NotifyError(r.Context(), "You do not have permission to delete users")
			goto response
		}

		if err := userSys.DeleteUser(username); err != nil {
			webutils.LogError(r.Context(), "Failed to delete user", err)
			webutils.NotifyError(r.Context(), "Could not delete user")
			htmx.Reswap(w, "none")
			goto response
		}

		htmx.Reswap(w, "delete")
		htmx.Retarget(w, "tr[data-user="+strconv.Quote(username)+"]")
		webutils.NotifyInfo(r.Context(), "User deleted")

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
