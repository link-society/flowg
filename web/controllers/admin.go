package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/templates/components"
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

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
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

			notifications = append(notifications, "&#10060; Could not fetch roles")
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

					notifications = append(notifications, fmt.Sprintf("&#10060; Could not fetch role '%s'", roleName))
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

			notifications = append(notifications, "&#10060; Could not fetch users")
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

					notifications = append(notifications, fmt.Sprintf("&#10060; Could not fetch user '%s'", username))
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

	mux.HandleFunc("GET /web/admin/roles/new/{$}", func(w http.ResponseWriter, r *http.Request) {
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

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if permissions.CanEditACLs {
			w.Header().Add("HX-Trigger", "htmx-custom-modal-open")

			h := templ.Handler(components.RoleForm(components.RoleFormProps{
				Name: "",
				Scopes: []struct {
					Name     auth.Scope
					Selected bool
				}{
					{auth.SCOPE_READ_STREAMS, false},
					{auth.SCOPE_WRITE_STREAMS, false},
					{auth.SCOPE_READ_TRANSFORMERS, false},
					{auth.SCOPE_WRITE_TRANSFORMERS, false},
					{auth.SCOPE_READ_PIPELINES, false},
					{auth.SCOPE_WRITE_PIPELINES, false},
					{auth.SCOPE_READ_ACLS, false},
					{auth.SCOPE_WRITE_ACLS, false},
					{auth.SCOPE_SEND_LOGS, false},
				},
			}))
			h.ServeHTTP(w, r)
		} else {
			trigger := map[string]interface{}{
				"htmx-custom-modal-open": map[string]interface{}{},
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

				triggerData = []byte("htmx-custom-modal-open")
			}

			w.Header().Add("HX-Trigger", string(triggerData))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to create roles"))
		}
	})

	mux.HandleFunc("POST /web/admin/roles/new/{$}", func(w http.ResponseWriter, r *http.Request) {
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

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if permissions.CanEditACLs {
			props := components.RoleFormProps{
				Name: "",
				Scopes: []struct {
					Name     auth.Scope
					Selected bool
				}{
					{auth.SCOPE_READ_STREAMS, false},
					{auth.SCOPE_WRITE_STREAMS, false},
					{auth.SCOPE_READ_TRANSFORMERS, false},
					{auth.SCOPE_WRITE_TRANSFORMERS, false},
					{auth.SCOPE_READ_PIPELINES, false},
					{auth.SCOPE_WRITE_PIPELINES, false},
					{auth.SCOPE_READ_ACLS, false},
					{auth.SCOPE_WRITE_ACLS, false},
					{auth.SCOPE_SEND_LOGS, false},
				},
			}

			trigger := map[string]interface{}{}

			err := r.ParseForm()
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error parsing form",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not parse form")

			} else {
				role := auth.Role{}

				props.Name = r.FormValue("name")
				role.Name = props.Name

				for i, scope := range props.Scopes {
					props.Scopes[i].Selected = r.FormValue(string(scope.Name)) == "on"

					if props.Scopes[i].Selected {
						role.Scopes = append(role.Scopes, scope.Name)
					}
				}

				err := authDb.SaveRole(role)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving role",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "&#10060; Could not save role")
				} else {
					trigger["htmx-custom-modal-close"] = map[string]interface{}{
						"after": "reload",
					}
				}
			}

			trigger["htmx-custom-toast"] = map[string]interface{}{
				"messages": notifications,
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

			h := templ.Handler(components.RoleForm(props))
			h.ServeHTTP(w, r)
		} else {
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

				triggerData = []byte("htmx-custom-modal-open")
			}

			w.Header().Add("HX-Trigger", string(triggerData))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to create roles"))
		}
	})

	mux.HandleFunc("GET /web/admin/roles/delete/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
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

			notifications = append(notifications, "&#10060; Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanEditACLs {
			notifications = append(notifications, "&#10060; You do not have permission to create roles")
		} else {
			roleName := r.PathValue("name")

			err := authDb.DeleteRole(roleName)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error deleting role",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not delete role")
			} else {
				w.Header().Add("HX-Reswap", "delete")
				w.Header().Add("HX-Retarget", "tr[data-role="+strconv.Quote(roleName)+"]")

				notifications = append(notifications, "&#9989; Role deleted")
			}
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
	})

	return mux
}
