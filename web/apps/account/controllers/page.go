package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"

	"link-society.com/flowg/web/apps/account/templates/views"
)

func Page(
	userSys *auth.UserSystem,
	tokenSys *auth.TokenSystem,
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

		tokenUUIDs, err := tokenSys.ListTokens(user.Name)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing personal access tokens",
				"channel", "web",
				"user", user.Name,
				"error", err.Error(),
			)

			notifications = append(notifications, "&#10060; Could not fetch personal access tokens")
			tokenUUIDs = []string{}
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				User:       user,
				TokenUUIDs: tokenUUIDs,

				Permissions:   permissions,
				Notifications: notifications,
			},
		))
		h.ServeHTTP(w, r)
	}
}
