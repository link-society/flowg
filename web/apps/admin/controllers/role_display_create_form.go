package controllers

import (
	"log/slog"

	"encoding/json"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func DisplayRoleCreateForm(
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
	}
}
