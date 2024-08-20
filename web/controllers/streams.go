package controllers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"

	"link-society.com/flowg/internal/filterdsl"
	"link-society.com/flowg/internal/storage"

	"link-society.com/flowg/web/templates/views"
)

func StreamController(db *storage.Storage) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /web/streams/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		streams, err := db.ListStreams()
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
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /web/streams/{name}/{$}", func(w http.ResponseWriter, r *http.Request) {
		notifications := []string{}

		streams, err := db.ListStreams()
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
		fields, err := db.ListStreamFields(stream)
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

		toRaw := r.URL.Query().Get("to")
		var to time.Time
		if toRaw != "" {
			to, err = time.Parse(time.RFC3339, toRaw)
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
			to = time.Now()
		}

		fromRaw := r.URL.Query().Get("from")
		var from time.Time
		if fromRaw != "" {
			from, err = time.Parse(time.RFC3339, fromRaw)
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
			from = to.Add(-5 * time.Minute)
		}

		var filter storage.Filter

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

		logs, err := db.Query(r.Context(), stream, from, to, filter)
		if err != nil {
			slog.ErrorContext(
				r.Context(),
				"error querying logs",
				"channel", "web",
				"stream", stream,
				"from", from,
				"to", to,
				"filter", filter,
				"error", err.Error(),
			)
			logs = []storage.LogEntry{}
			notifications = append(notifications, "❌ Could not fetch logs")
		}

		h := templ.Handler(views.Streams(
			views.StreamsProps{
				Streams:       streams,
				CurrentStream: stream,

				LogEntries: logs,
				Fields:     fields,
				From:       from,
				To:         to,
				Filter:     filterSource,
			},
			notifications,
		))
		h.ServeHTTP(w, r)
	})

	return mux
}
