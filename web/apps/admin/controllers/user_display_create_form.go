package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func DisplayUserCreateForm(
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
			roleFields := []struct {
				Name     string
				Selected bool
			}{}

			roles, err := roleSys.ListRoles()
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error listing roles",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not fetch roles")
			} else {
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
			}

			trigger := htmx.Trigger{
				ModalOpenEvent: &htmx.ModalOpenEvent{},
				ToastEvent: &htmx.ToastEvent{
					Messages: notifications,
				},
			}

			trigger.Write(r.Context(), w)
			h := templ.Handler(components.UserForm(components.UserFormProps{
				Name:     "",
				Password: "",
				Roles:    roleFields,
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
			w.Write([]byte("&#10060; You do not have permission to create users"))
		}
	}
}
