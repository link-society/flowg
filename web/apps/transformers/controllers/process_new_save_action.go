package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"

	"link-society.com/flowg/web/apps/transformers/templates/views"
)

func ProcessNewSaveAction(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
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

		if !permissions.CanViewTransformers {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		transformerName := ""
		transformerCode := "."

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
			transformerName = r.FormValue("name")
			transformerCode = r.FormValue("code")

			if permissions.CanEditTransformers {
				if transformerName == "" {
					notifications = append(notifications, "&#10060; Transformer name is required")
				}

				if transformerCode == "" {
					notifications = append(notifications, "&#10060; Transformer code is required")
				}

				if transformerName != "" && transformerCode != "" {
					err = pipelinesManager.SaveTransformerScript(transformerName, transformerCode)
					if err != nil {
						slog.ErrorContext(
							r.Context(),
							"error saving transformer script",
							"channel", "web",
							"error", err.Error(),
						)

						notifications = append(notifications, "&#10060; Could not save transformer script")
					} else {
						notifications = append(notifications, "&#9989; Transformer script saved")
					}
				}
			}
		}

		transformers, err := pipelinesManager.ListTransformers()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing transformers",
				"channel", "web",
				"error", err.Error(),
			)

			transformers = []string{}
			notifications = append(notifications, "&#10060; Could not fetch transformers")
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Transformers:       transformers,
				CurrentTransformer: transformerName,
				Code:               transformerCode,

				Permissions:   permissions,
				Notifications: notifications,
			},
		))
		h.ServeHTTP(w, r)
	}
}
