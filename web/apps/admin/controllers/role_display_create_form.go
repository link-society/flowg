package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils/htmx"

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
			trigger := htmx.Trigger{
				ModalOpenEvent: &htmx.ModalOpenEvent{},
			}
			trigger.Write(r.Context(), w)

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
			trigger := htmx.Trigger{
				ModalOpenEvent: &htmx.ModalOpenEvent{},
				ToastEvent: &htmx.ToastEvent{
					Messages: notifications,
				},
			}

			trigger.Write(r.Context(), w)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to create roles"))
		}
	}
}
