package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/admin/templates/views"
)

func Page(
	roleSys *auth.RoleSystem,
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewACLs {
			http.Redirect(w, r, "/web/", http.StatusSeeOther)
			return
		}

		roles, err := roleSys.ListRoles()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch roles", err)
			webutils.NotifyError(r.Context(), "Could not fetch roles")
			roles = []auth.Role{}
		}

		users, err := userSys.ListUsers()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch users", err)
			webutils.NotifyError(r.Context(), "Could not fetch users")
			users = []auth.User{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Roles: roles,
				Users: users,
			},
		))
		h.ServeHTTP(w, r)
	}
}
