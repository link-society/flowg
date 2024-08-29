package controllers

import (
	"log/slog"

	"encoding/json"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func ProcessRoleCreateForm(
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

				err := roleSys.SaveRole(role)
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
	}
}
