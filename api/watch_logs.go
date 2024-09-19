package api

import (
	"context"
	"log/slog"

	"encoding/json"
	"fmt"

	"net/http"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/internal/ffi/filterdsl"
)

type WatchLogsRequest struct {
	Stream string  `path:"stream" minLength:"1"`
	Filter *string `query:"filter"`
}

type WatchLogsResponse struct {
	usecase.OutputWithEmbeddedWriter
}

func WatchLogsUsecase(
	authDb *auth.Database,
	logNotifier *lognotify.LogNotifier,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req WatchLogsRequest,
				resp *WatchLogsResponse,
			) error {
				var filter logstorage.Filter

				if req.Filter != nil {
					var err error

					filter, err = filterdsl.Compile(*req.Filter)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to compile filter",
							"channel", "api",
							"error", err.Error(),
							"filter", req.Filter,
						)

						return status.Wrap(err, status.InvalidArgument)
					}
				}

				resp.Writer.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")
				resp.Writer.(http.ResponseWriter).Header().Set("Access-Control-Expose-Headers", "Content-Type")

				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
				resp.Writer.(http.ResponseWriter).Header().Set("Connection", "keep-alive")

				slog.InfoContext(
					ctx,
					"watch logs",
					"channel", "api",
					"stream", req.Stream,
				)
				defer slog.InfoContext(
					ctx,
					"done watching logs",
					"channel", "api",
					"stream", req.Stream,
				)

				logC := logNotifier.Subscribe(req.Stream, ctx.Done())

				for log := range logC {
					if filter == nil || filter.Evaluate(&log.LogEntry) {
						payload, err := json.Marshal(log.LogEntry)
						if err != nil {
							fmt.Fprintf(resp, "event: error\n")
							fmt.Fprintf(resp, "data: %s\n\n", err.Error())
							resp.Writer.(http.Flusher).Flush()

							return nil
						}

						fmt.Fprintf(resp, "id: %s\n", log.LogKey)
						fmt.Fprintf(resp, "event: log\n")
						fmt.Fprintf(resp, "data: %s\n\n", payload)
						resp.Writer.(http.Flusher).Flush()
					}
				}

				return nil
			},
		),
	)

	u.SetName("watch_logs")
	u.SetTitle("Watch Logs")
	u.SetDescription("Server-Sent Events endpoint to watch logs in real time.")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
