package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
	"link-society.com/flowg/internal/utils/langs/filtering"
)

// QueryStreamDeps lists the dependencies of [NewQueryStreamUsecase].
type QueryStreamDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewQueryStreamUsecase retrieves the logs of a stream within a time range,
// optionally narrowed by indexes and a filter expression.
//
// Callers must have the read-streams permission. A filter that fails to compile
// is reported as an invalid-argument error.
func NewQueryStreamUsecase(deps QueryStreamDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req schemas.QueryStreamRequest,
				resp *schemas.QueryStreamResponse,
			) error {
				var filter filtering.Filter

				if req.Filter != nil {
					var err error
					filter, err = filtering.Compile(*req.Filter)
					if err != nil {
						logger.ErrorContext(
							ctx,
							"Failed to compile filter",
							slog.String("stream", req.Stream),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						resp.Records = nil
						return status.Wrap(err, status.InvalidArgument)
					}
				} else {
					filter = nil
				}

				records, err := deps.LogStorage.FetchLogs(ctx, req.Stream, req.From, req.To, filter, req.Indexing)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to query logs",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Records = records
				return nil
			},
		),
	)

	u.SetName("query_stream")
	u.SetTitle("Query Stream")
	u.SetDescription("Query logs from a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewQueryStreamUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}/logs",
	)
}
