package controllers

import (
	"strconv"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessRoleDeleteAction(
	roleSys *auth.RoleSystem,
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		roleName := r.PathValue("name")

		if !webutils.Permissions(r.Context()).CanEditACLs {
			htmx.Reswap(w, "none")
			webutils.NotifyError(r.Context(), "You do not have permission to delete roles")
			goto response
		}

		if err := roleSys.DeleteRole(roleName); err != nil {
			webutils.LogError(r.Context(), "Failed to delete role", err)
			webutils.NotifyError(r.Context(), "Could not delete role")
			htmx.Reswap(w, "none")
			goto response
		}

		htmx.Reswap(w, "delete")
		htmx.Retarget(w, "tr[data-role="+strconv.Quote(roleName)+"]")
		webutils.NotifyInfo(r.Context(), "Role deleted")

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
