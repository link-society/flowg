package controllers

import (
	"net/http"
	"sort"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"

	"link-society.com/flowg/web/apps/storage/templates/views"
)

func StreamPage(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewStreams {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		streamNames := []string{}
		streams, err := metaSys.ListStreams()
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch streams", err)
			webutils.NotifyError(r.Context(), "Could not fetch streams")
		} else {
			for streamName := range streams {
				streamNames = append(streamNames, streamName)
			}
		}

		streamName := r.PathValue("name")
		streamConfig, err := metaSys.GetStreamConfig(streamName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch stream config", err)
			webutils.NotifyError(r.Context(), "Could not fetch stream config")
		}

		sort.Strings(streamNames)

		h := templ.Handler(views.Page(
			views.PageProps{
				StreamNames:         streamNames,
				CurrentStreamName:   streamName,
				CurrentStreamConfig: &streamConfig,

				Permissions:   webutils.Permissions(r.Context()),
				Notifications: webutils.Notifications(r.Context()),
			},
		))
		h.ServeHTTP(w, r)
	}
}
