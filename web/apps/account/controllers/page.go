package controllers

import (
	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/account/templates/views"
)

func Page(
	userSys *auth.UserSystem,
	tokenSys *auth.TokenSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))
		user := auth.GetContextUser(r.Context())

		tokenUUIDs, err := tokenSys.ListTokens(user.Name)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to list personal access tokens", err)
			webutils.NotifyError(r.Context(), "Could not fetch personal access tokens")
			tokenUUIDs = []string{}
			goto response
		}

	response:
		h := templ.Handler(views.Page(
			views.PageProps{
				User:       user,
				TokenUUIDs: tokenUUIDs,
			},
		))
		h.ServeHTTP(w, r)
	}
}
