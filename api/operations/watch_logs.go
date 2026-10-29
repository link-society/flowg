package operations

import (
	"context"
	"log/slog"

	"encoding/json"
	"fmt"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	"go.link-society.com/flowg/internal/engines/lognotify"
	"go.link-society.com/flowg/internal/models"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/internal/utils/langs/filtering"
)

// WatchLogsDeps lists the dependencies of [NewWatchLogsUsecase].
type WatchLogsDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage storage.AuthStorage
	LogNotifier lognotify.LogNotifier
}

// NewWatchLogsUsecase streams a stream's logs to the client in real time as
// Server-Sent Events.
//
// It subscribes to new records and forwards those matching the optional filter
// until the client disconnects. Callers must have the read-streams permission.
// A filter that fails to compile is reported as an invalid-argument error.
func NewWatchLogsUsecase(deps WatchLogsDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req schemas.WatchLogsRequest,
				resp *schemas.WatchLogsResponse,
			) error {
				var filter filtering.Filter

				if req.Filter != nil && *req.Filter != "" {
					var err error

					filter, err = filtering.Compile(*req.Filter)
					if err != nil {
						deps.Logger.ErrorContext(
							ctx,
							"Failed to compile filter",
							slog.String("error", err.Error()),
							slog.String("stream", req.Stream),
							slog.String("filter", *req.Filter),
						)

						return status.Wrap(err, status.InvalidArgument)
					}
				}

				logM, err := deps.LogNotifier.Subscribe(ctx, req.Stream)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to subscribe to logs",
						slog.String("error", err.Error()),
						slog.String("stream", req.Stream),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Writer.(http.ResponseWriter).Header().Set("Access-Control-Allow-Origin", "*")
				resp.Writer.(http.ResponseWriter).Header().Set("Access-Control-Expose-Headers", "Content-Type")

				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "text/event-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")
				resp.Writer.(http.ResponseWriter).Header().Set("Connection", "keep-alive")

				deps.Logger.InfoContext(
					ctx,
					"watch logs",
					slog.String("stream", req.Stream),
				)
				defer deps.Logger.InfoContext(
					ctx,
					"done watching logs",
					slog.String("stream", req.Stream),
				)

				for log := range logM.ReceiveC() {
					matches := false

					if filter == nil {
						matches = true
					} else {
						var err error
						matches, err = filter.Evaluate(&log.LogRecord)
						if err != nil {
							fmt.Fprintf(resp.Writer, "event: error\n")
							fmt.Fprintf(
								resp.Writer,
								"data: failed to evaluate filter for log entry '%s': %s\n\n",
								log.LogKey,
								err.Error(),
							)
							resp.Writer.(http.Flusher).Flush()

							return nil
						}
					}

					if matches {
						payload, err := json.Marshal(log.LogRecord)
						if err != nil {
							fmt.Fprintf(resp.Writer, "event: error\n")
							fmt.Fprintf(resp.Writer, "data: %s\n\n", err.Error())
							resp.Writer.(http.Flusher).Flush()

							return nil
						}

						fmt.Fprintf(resp.Writer, "id: %s\n", log.LogKey)
						fmt.Fprintf(resp.Writer, "event: log\n")
						fmt.Fprintf(resp.Writer, "data: %s\n\n", payload)
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

// annotateWatchLogs documents the watch response as a Server-Sent Events stream.
func annotateWatchLogs(oc openapi.OperationContext) error {
	contentUnits := oc.Response()
	for i, cu := range contentUnits {
		if cu.HTTPStatus == 200 {
			cu.ContentType = "text/event-stream"
			cu.Description = "Stream of log entries"
			cu.Format = "Server-Sent Events"
		}

		contentUnits[i] = cu
	}

	return nil
}

func init() {
	routing.RegisterOperation(
		NewWatchLogsUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}/logs/watch",
		routing.Annotated(annotateWatchLogs),
	)
}
