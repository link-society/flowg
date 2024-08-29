package controllers

import (
	"fmt"
	"log/slog"

	"encoding/json"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"

	"link-society.com/flowg/web/apps/admin/templates/components"
)

func ProcessUserCreateForm(
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

			props := components.UserFormProps{
				Name:     "",
				Password: "",
				Roles:    roleFields,
			}

			trigger := map[string]interface{}{}

			err = r.ParseForm()
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error parsing form",
					"channel", "web",
					"error", err.Error(),
				)

				notifications = append(notifications, "&#10060; Could not parse form")

			} else {
				user := auth.User{}

				props.Name = r.FormValue("name")
				props.Password = r.FormValue("password")
				user.Name = props.Name

				for i, role := range props.Roles {
					props.Roles[i].Selected = r.FormValue(fmt.Sprintf("role_%s", role.Name)) == "on"

					if props.Roles[i].Selected {
						user.Roles = append(user.Roles, role.Name)
					}
				}

				err := userSys.SaveUser(user, props.Password)
				if err != nil {
					slog.ErrorContext(
						r.Context(),
						"error saving user",
						"channel", "web",
						"error", err.Error(),
					)

					notifications = append(notifications, "&#10060; Could not save user")
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

			h := templ.Handler(components.UserForm(props))
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
			w.Write([]byte("&#10060; You do not have permission to create users"))
		}
	}
}
