package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/ffi/filterdsl"
	"link-society.com/flowg/internal/webutils"
	"link-society.com/flowg/internal/webutils/htmx"

	"link-society.com/flowg/web/apps/streams/templates/components"
	"link-society.com/flowg/web/apps/streams/templates/views"
)

func StreamPage(
	userSys *auth.UserSystem,
	metaSys *logstorage.MetaSystem,
	querySys *logstorage.QuerySystem,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(webutils.WithNotificationSystem(r.Context()))
		r = r.WithContext(webutils.WithPermissionSystem(r.Context(), userSys))

		if !webutils.Permissions(r.Context()).CanViewStreams {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		// fetch data for template
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

		stream := r.PathValue("name")
		fields, err := metaSys.ListStreamFields(stream)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch fields", err)
			webutils.NotifyError(r.Context(), "Could not fetch fields")
			fields = []string{}
		}

		// parse filter values in querystring
		autoRefresh := r.URL.Query().Get("autorefresh")
		if autoRefresh == "" {
			autoRefresh = "0"
		}

		toRaw := r.URL.Query().Get("to")
		var naiveTo time.Time
		if toRaw != "" {
			naiveTo, err = time.Parse("2006-01-02T15:04:05", toRaw)
			if err != nil {
				webutils.LogError(r.Context(), "Failed to parse 'to' time", err)
				webutils.NotifyError(r.Context(), "Invalid 'to' time")
			}
		} else {
			naiveTo = time.Now()
		}

		fromRaw := r.URL.Query().Get("from")
		var naiveFrom time.Time
		if fromRaw != "" {
			naiveFrom, err = time.Parse("2006-01-02T15:04:05", fromRaw)
			if err != nil {
				webutils.LogError(r.Context(), "Failed to parse 'from' time", err)
				webutils.NotifyError(r.Context(), "Invalid 'from' time")
			}
		} else {
			naiveFrom = naiveTo.Add(-5 * time.Minute)
		}

		timeOffsetRaw := r.URL.Query().Get("timeoffset")
		var timeOffset int
		if timeOffsetRaw != "" {
			timeOffset, err = strconv.Atoi(timeOffsetRaw)
			if err != nil {
				webutils.LogError(r.Context(), "Failed to parse 'timeoffset'", err)
				webutils.NotifyError(r.Context(), "Invalid 'timeoffset'")
			}
		} else {
			timeOffset = 0
		}

		var filter logstorage.Filter

		filterSource := r.URL.Query().Get("filter")
		if filterSource != "" {
			filter, err = filterdsl.Compile(filterSource)
			if err != nil {
				webutils.LogError(r.Context(), "Failed to compile filter", err)
				webutils.NotifyError(r.Context(), "Invalid filter")
			}
		} else {
			filter = nil
		}

		// fetch logs
		localFrom := naiveFrom.Add(time.Duration(timeOffset) * time.Minute)
		localTo := naiveTo.Add(time.Duration(timeOffset) * time.Minute)

		logs, err := querySys.FetchLogs(r.Context(), stream, localFrom, localTo, filter)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to fetch logs", err)
			webutils.NotifyError(r.Context(), "Could not fetch logs")
			logs = []logstorage.LogEntry{}
		}

		// aggregate for histogram
		interval := localTo.Sub(localFrom) / 24
		counts := make([]int, 24)
		timestamps := make([]time.Time, 24)

		for i := 0; i < 24; i++ {
			start := localFrom.Add(time.Duration(i) * interval)
			end := localFrom.Add(time.Duration(i+1) * interval)

			timestamps[i] = start

			for _, log := range logs {
				if log.Timestamp.After(start) && log.Timestamp.Before(end) {
					counts[i]++
				}
			}
		}

		var timeserie [][2]interface{}
		for i := 0; i < 24; i++ {
			timeserie = append(timeserie, [2]interface{}{
				timestamps[i].UnixMilli(),
				counts[i],
			})
		}

		histogramData, err := json.Marshal(timeserie)
		if err != nil {
			webutils.LogError(r.Context(), "Failed to marshal histogram data", err)
			webutils.NotifyError(r.Context(), "Could not fetch histogram data")
			histogramData = []byte("[]")
		}

		// render
		if htmx.IsHtmxRequest(r) {
			trigger := htmx.Trigger{
				ToastEvent: &htmx.ToastEvent{
					Messages: webutils.Notifications(r.Context()),
				},
			}

			trigger.Write(r.Context(), w)
			h := templ.Handler(components.Viewer(
				components.ViewerProps{
					LogEntries:    logs,
					Fields:        fields,
					From:          naiveFrom,
					To:            naiveTo,
					Filter:        filterSource,
					AutoRefresh:   autoRefresh,
					HistogramData: string(histogramData),
				},
			))
			h.ServeHTTP(w, r)
		} else {
			h := templ.Handler(views.Page(
				views.PageProps{
					Streams:       streamNames,
					CurrentStream: stream,

					LogEntries:  logs,
					Fields:      fields,
					From:        naiveFrom,
					To:          naiveTo,
					Filter:      filterSource,
					AutoRefresh: autoRefresh,

					HistogramData: string(histogramData),

					Permissions:   webutils.Permissions(r.Context()),
					Notifications: webutils.Notifications(r.Context()),
				},
			))
			h.ServeHTTP(w, r)
		}
	}
}
