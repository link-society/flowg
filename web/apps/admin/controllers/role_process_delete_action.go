package controllers

import (
	"log/slog"

	"strconv"

	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessRoleDeleteAction(
	roleSys *auth.RoleSystem,
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := userSys.ListUserScopes(user.Name)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanEditACLs {
			htmx.Reswap(w, "none")

			notifications = append(notifications, "&#10060; You do not have permission to delete roles")
		} else {
			roleName := r.PathValue("name")

			err := roleSys.DeleteRole(roleName)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error deleting role",
					"channel", "web",
					"error", err.Error(),
				)

				htmx.Reswap(w, "none")

				notifications = append(notifications, "&#10060; Could not delete role")
			} else {
				htmx.Reswap(w, "delete")
				htmx.Retarget(w, "tr[data-role="+strconv.Quote(roleName)+"]")

				notifications = append(notifications, "&#9989; Role deleted")
			}
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
