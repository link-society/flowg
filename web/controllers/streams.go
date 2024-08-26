package controllers

import (
	"log/slog"
	"time"

	"encoding/json"
	"strconv"

	"net/http"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/filterdsl"
	"link-society.com/flowg/internal/logstorage"

	"link-society.com/flowg/web/templates/components"
	"link-society.com/flowg/web/templates/views"
)

func StreamController(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/streams/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
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
			notifications = append(notifications, "❌ Could not fetch streams")
		}

		if len(streams) > 0 {
			defaultStream := streams[0]

			http.Redirect(w, r, "/web/streams/"+defaultStream+"/", http.StatusFound)
			return
		}

		h := templ.Handler(views.Streams(
			views.StreamsProps{
				Streams:       streams,
				CurrentStream: "",
			},
			permissions,
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/streams/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		permissions := auth.Permissions{}
		notifications := []string{}

		user := auth.GetContextUser(r.Context())
		scopes, err := authDb.ListUserScopes(user)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing user scopes",
				"channel", "web",
				"error", err.Error(),
			)

			notifications = append(notifications, "❌ Could not fetch user permissions")
		} else {
			permissions = auth.PermissionsFromScopes(scopes)
		}

		if !permissions.CanViewStreams {
			http.Redirect(w, r, "/web", http.StatusSeeOther)
			return
		}

		// fetch data for template
		streams, err := logDb.ListStreams()
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing streams",
				"channel", "web",
				"error", err.Error(),
			)

			streams = []string{}
			notifications = append(notifications, "❌ Could not fetch streams")
		}

		stream := r.PathValue("name")
		fields, err := logDb.ListStreamFields(stream)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error listing fields",
				"channel", "web",
				"stream", stream,
				"error", err.Error(),
			)

			fields = []string{}
			notifications = append(notifications, "❌ Could not fetch fields")
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
				slog.ErrorContext(
					r.Context(),
					"error parsing 'to' time",
					"channel", "web",
					"stream", stream,
					"to", toRaw,
					"error", err.Error(),
				)
				notifications = append(notifications, "❌ Invalid 'to' time")
			}
		} else {
			naiveTo = time.Now()
		}

		fromRaw := r.URL.Query().Get("from")
		var naiveFrom time.Time
		if fromRaw != "" {
			naiveFrom, err = time.Parse("2006-01-02T15:04:05", fromRaw)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error parsing 'from' time",
					"channel", "web",
					"stream", stream,
					"from", fromRaw,
					"error", err.Error(),
				)
				notifications = append(notifications, "❌ Invalid 'from' time")
			}
		} else {
			naiveFrom = naiveTo.Add(-5 * time.Minute)
		}

		timeOffsetRaw := r.URL.Query().Get("timeoffset")
		var timeOffset int
		if timeOffsetRaw != "" {
			timeOffset, err = strconv.Atoi(timeOffsetRaw)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error parsing 'timeoffset'",
					"channel", "web",
					"stream", stream,
					"timeoffset", timeOffsetRaw,
					"error", err.Error(),
				)
				notifications = append(notifications, "❌ Invalid 'timeoffset'")
			}
		} else {
			timeOffset = 0
		}

		var filter logstorage.Filter

		filterSource := r.URL.Query().Get("filter")
		if filterSource != "" {
			filter, err = filterdsl.Compile(filterSource)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"error compiling filter",
					"channel", "web",
					"stream", stream,
					"filter", filterSource,
					"error", err.Error(),
				)
				notifications = append(notifications, "❌ Invalid filter")
			}
		} else {
			filter = nil
		}

		// fetch logs
		localFrom := naiveFrom.Add(time.Duration(timeOffset) * time.Minute)
		localTo := naiveTo.Add(time.Duration(timeOffset) * time.Minute)

		logs, err := logDb.Query(r.Context(), stream, localFrom, localTo, filter)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error querying logs",
				"channel", "web",
				"stream", stream,
				"from", localFrom,
				"to", localTo,
				"filter", filter,
				"error", err.Error(),
			)
			logs = []logstorage.LogEntry{}
			notifications = append(notifications, "❌ Could not fetch logs")
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
			slog.ErrorContext(
				r.Context(),
				"error marshalling histogram data",
				"channel", "web",
				"stream", stream,
				"error", err.Error(),
			)

			histogramData = []byte("[]")
			notifications = append(notifications, "❌ Could not fetch histogram data")
		}

		// render
		if r.Header.Get("HX-Request") == "true" {
			h := templ.Handler(components.StreamViewer(
				components.StreamViewerProps{
					LogEntries:    logs,
					Fields:        fields,
					From:          naiveFrom,
					To:            naiveTo,
					Filter:        filterSource,
					AutoRefresh:   autoRefresh,
					HistogramData: string(histogramData),

					Notifications: notifications,
				},
			))
			h.ServeHTTP(w, r)
		} else {
			h := templ.Handler(views.Streams(
				views.StreamsProps{
					Streams:       streams,
					CurrentStream: stream,

					LogEntries:  logs,
					Fields:      fields,
					From:        naiveFrom,
					To:          naiveTo,
					Filter:      filterSource,
					AutoRefresh: autoRefresh,

					HistogramData: string(histogramData),
				},
				permissions,
				notifications,
			))
			h.ServeHTTP(w, r)
		}
	})

	return mux
}
