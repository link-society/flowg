package api

import (
	"context"
	"log/slog"

	"encoding/json"
	"fmt"

	"net/http"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"
	"link-society.com/flowg/internal/utils/ffi/filterdsl"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/storage/auth"
)

type WatchLogsRequest struct {
	Stream string  `path:"stream" minLength:"1"`
	Filter *string `query:"filter"`
}

type WatchLogsResponse struct {
	usecase.OutputWithEmbeddedWriter
}

func WatchLogsUsecase(
	authStorage *auth.Storage,
	logNotifier *lognotify.LogNotifier,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req WatchLogsRequest,
				resp *WatchLogsResponse,
			) error {
				var filter filterdsl.Filter

				if req.Filter != nil && *req.Filter != "" {
					var err error

					filter, err = filterdsl.Compile(*req.Filter)
					if err != nil {
						slog.ErrorContext(
							ctx,
							"Failed to compile filter",
							slog.String("channel", "api"),
							slog.String("error", err.Error()),
							slog.String("stream", req.Stream),
							slog.String("filter", *req.Filter),
						)

						return status.Wrap(err, status.InvalidArgument)
					}
				}

				logM, err := logNotifier.Subscribe(ctx, req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to subscribe to logs",
						slog.String("channel", "api"),
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

				slog.InfoContext(
					ctx,
					"watch logs",
					slog.String("channel", "api"),
					slog.String("stream", req.Stream),
				)
				defer slog.InfoContext(
					ctx,
					"done watching logs",
					slog.String("channel", "api"),
					slog.String("stream", req.Stream),
				)

				for log := range logM.ReceiveC() {
					if filter == nil || filter.Evaluate(&log.LogRecord) {
						payload, err := json.Marshal(log.LogRecord)
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
