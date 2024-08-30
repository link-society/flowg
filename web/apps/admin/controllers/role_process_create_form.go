package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func ProcessRoleCreateForm(
	roleSys *auth.RoleSystem,
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanEditACLs {
			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to create roles"))
			return
		}

		var (
			trigger htmx.Trigger
			role    auth.Role
		)

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

		err := r.ParseForm()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		props.Name = r.FormValue("name")
		role.Name = props.Name

		for i, scope := range props.Scopes {
			props.Scopes[i].Selected = r.FormValue(string(scope.Name)) == "on"

			if props.Scopes[i].Selected {
				role.Scopes = append(role.Scopes, scope.Name)
			}
		}

		err = roleSys.SaveRole(role)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to save role", err)
			webutils.NotifyError(r.Context(), "Could not save role")
			goto response
		}

		trigger.ModalCloseEvent = &htmx.ModalCloseEvent{
			After: "reload",
		}

	response:
		trigger.ToastEvent = &htmx.ToastEvent{
			Messages: webutils.Notifications(r.Context()),
		}

		trigger.Write(r.Context(), w)
		h := templ.Handler(components.RoleForm(props))
		h.ServeHTTP(w, r)
	}
}
