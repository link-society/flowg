package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"

	"link-society.com/flowg/web/apps/transformers/templates/views"
)

func PageEdit(
	userSys *auth.UserSystem,
	pipelinesManager *pipelines.Manager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transformerName := r.PathValue("name")
		transformerCode, err := pipelinesManager.GetTransformerScript(transformerName)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error getting transformer script",
				"channel", "web",
				"error", err.Error(),
			)
			http.Redirect(w, r, "/web/transformers/new", http.StatusTemporaryRedirect)
			return
		}

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
