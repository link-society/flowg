package controllers

import (
	"log/slog"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/logstorage"

	"link-society.com/flowg/web/apps/streams/templates/views"
)

func DefaultPage(
	userSys *auth.UserSystem,
	logDb *logstorage.Storage,
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

		if !permissions.CanViewStreams {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}
		streams, err := logDb.ListStreams()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing streams",
				"channel", "web",
				"error", err.Error(),
			)

			streams = []string{}
			notifications = append(notifications, "&#10060; Could not fetch streams")
		}

		if len(streams) > 0 {
			defaultStream := streams[0]

			http.Redirect(w, r, "/web/streams/"+defaultStream+"/", http.StatusFound)
			return
		}

		h := templ.Handler(views.Page(
			views.PageProps{
				Streams:       streams,
				CurrentStream: "",

				Permissions:   permissions,
				Notifications: notifications,
			},
		))
		h.ServeHTTP(w, r)
	}
}
