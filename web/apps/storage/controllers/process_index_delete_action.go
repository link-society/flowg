package controllers

import (
	"net/http"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessIndexDeleteAction(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanEditStreams {
			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "none")
			w.WriteHeader(http.StatusOK)
			return
		}

		streamName := r.PathValue("name")
		fieldName := r.PathValue("field")

		streamConfig, err := metaSys.GetStreamConfig(streamName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch stream config", err)
			webutils.NotifyError(r.Context(), "Failed to fetch stream config")

			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "none")
			w.WriteHeader(http.StatusOK)
			return
		}

		for i, f := range streamConfig.IndexedFields {
			if f == fieldName {
				streamConfig.IndexedFields = append(
					streamConfig.IndexedFields[:i],
					streamConfig.IndexedFields[i+1:]...,
				)
				break
			}
		}

		if err := metaSys.ConfigureStream(streamName, streamConfig); err != nil {
			webutils.LogError(r.Context(), "Failed to save stream config", err)
			webutils.NotifyError(r.Context(), "Failed to save stream config")

			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			htmx.Reswap(w, "none")
			w.WriteHeader(http.StatusOK)
			return
		}

		webutils.NotifyInfo(r.Context(), "Field unindexed")

		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}
		trigger.Write(r.Context(), w)
		htmx.Reswap(w, "delete")
		w.WriteHeader(http.StatusOK)
	}
}
