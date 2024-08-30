package controllers

import (
	"fmt"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func ProcessUserCreateForm(
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

		props := components.UserFormProps{
			Name:     "",
			Password: "",
			Roles:    roleFields,
		}
		user := auth.User{}
		trigger := htmx.Trigger{}

		err = r.ParseForm()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		props.Name = r.FormValue("name")
		props.Password = r.FormValue("password")
		user.Name = props.Name

		for i, role := range props.Roles {
			props.Roles[i].Selected = r.FormValue(fmt.Sprintf("role_%s", role.Name)) == "on"

			if props.Roles[i].Selected {
				user.Roles = append(user.Roles, role.Name)
			}
		}

		err = userSys.SaveUser(user, props.Password)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to save user", err)
			webutils.NotifyError(r.Context(), "Could not save user")
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
		h := templ.Handler(components.UserForm(props))
		h.ServeHTTP(w, r)
	}
}
