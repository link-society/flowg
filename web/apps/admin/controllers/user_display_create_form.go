package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func DisplayUserCreateForm(
	roleSys *auth.RoleSystem,
	userSys *auth.UserSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanEditACLs {
			trigger := htmx.Trigger{
				ModalOpenEvent: &htmx.ModalOpenEvent{},
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("&#10060; You do not have permission to create users"))
			return
		}

		roles, err := roleSys.ListRoles()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch roles", err)
			webutils.NotifyError(r.Context(), "Could not fetch roles")
			roles = []auth.Role{}
		}

		roleFields := []struct {
			Name     string
			Selected bool
		}{}

		for _, role := range roles {
			roleFields = append(
				roleFields,
				struct {
					Name     string
					Selected bool
				}{
					Name:     role.Name,
					Selected: false,
				},
			)
		}

		trigger := htmx.Trigger{
			ModalOpenEvent: &htmx.ModalOpenEvent{},
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}

		trigger.Write(r.Context(), w)
		h := templ.Handler(components.UserForm(components.UserFormProps{
			Name:     "",
			Password: "",
			Roles:    roleFields,
		}))
		h.ServeHTTP(w, r)
	}
}
