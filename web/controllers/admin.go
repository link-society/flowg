package controllers

import (
	"fmt"
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/templates/views"
)

func AdminController(authDb *auth.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/admin/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewACLs {
			http.Redirect(w, r, "/web/", http.StatusSeeOther)
			return
		}

		roles := []auth.Role{}
		roleNames, err := authDb.ListRoles()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing roles",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch roles")
		} else {
			for _, roleName := range roleNames {
				role, err := authDb.GetRole(roleName)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error fetching role",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, fmt.Sprintf("❌ Could not fetch role '%s'", roleName))
					continue
				}

				roles = append(roles, role)
			}
		}

		users := []*auth.User{}
		usernames, err := authDb.ListUsers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing users",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch users")
		} else {
			for _, username := range usernames {
				user, err := authDb.GetUser(username)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error fetching user",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, fmt.Sprintf("❌ Could not fetch user '%s'", username))
					continue
				}

				users = append(users, user)
			}
		}

		h := templ.Handler(views.Admin(
			views.AdminProps{
				Roles: roles,
				Users: users,
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	return mux
}
