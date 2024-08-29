package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/admin/templates/views"
)

func Page(
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

		if !permissions.CanViewACLs {
			http.Redirect(w, r, "/web/", http.StatusSeeOther)
			return
		}

		roles, err := roleSys.ListRoles()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing roles",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch roles")

			roles = []auth.Role{}
		}

		users, err := userSys.ListUsers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing users",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch users")

			users = []auth.User{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Roles: roles,
				Users: users,

				Permissions:   permissions,
				Notifications: notifications,
			},
		))
		h.ServeHTTP(w, r)
	}
}
