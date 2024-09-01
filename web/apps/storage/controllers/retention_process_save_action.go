package controllers

import (
	"net/http"
	"strconv"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"
)

func ProcessRetentionSaveAction(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
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
			w.Write([]byte("&#10060; You do not have permission to edit streams"))
			return
		}

		var (
			retentionTime int64
			retentionSize int64

			retentionTimeStr string
			retentionSizeStr string
		)

		streamName := r.PathValue("name")
		streamConfig, err := metaSys.GetStreamConfig(streamName)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch stream config", err)
			webutils.NotifyError(r.Context(), "Could not fetch stream config")
			goto response
		}

		if err := r.ParseForm(); err != nil {
			webutils.LogError(r.Context(), "Failed to parse form data", err)
			webutils.NotifyError(r.Context(), "Could not parse form")
			goto response
		}

		retentionTimeStr = r.FormValue("retention_time")
		retentionSizeStr = r.FormValue("retention_size")

		if retentionTimeStr == "" {
			retentionTimeStr = "0"
		}
		if retentionSizeStr == "" {
			retentionSizeStr = "0"
		}

		retentionTime, err = strconv.ParseInt(retentionTimeStr, 10, 64)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to parse retention time", err)
			webutils.NotifyError(r.Context(), "Could not parse retention time")
			retentionTime = -1
		}

		retentionSize, err = strconv.ParseInt(retentionSizeStr, 10, 64)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to parse retention size", err)
			webutils.NotifyError(r.Context(), "Could not parse retention size")
			retentionSize = -1
		}

		if retentionTime < 0 || retentionSize < 0 {
			goto response
		}

		streamConfig.RetentionTime = retentionTime
		streamConfig.RetentionSize = retentionSize

		err = metaSys.ConfigureStream(streamName, *streamConfig)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to configure stream", err)
			webutils.NotifyError(r.Context(), "Could not configure stream")
			goto response
		}

		webutils.NotifyInfo(r.Context(), "Stream configuration updated")
		htmx.Reswap(w, "none")

	response:
		trigger := htmx.Trigger{
			ToastEvent: &htmx.ToastEvent{
				Messages: webutils.Notifications(r.Context()),
			},
		}

		trigger.Write(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	}
}
