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
	"link-society.com/flowg/internal/models"

	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/log"
)

// GetStreamUsageDeps lists the dependencies of [NewGetStreamUsageUsecase].
type GetStreamUsageDeps struct {
	fx.In

	AuthStorage authStorage.Storage
	LogStorage  log.Storage
}

// GetStreamUsageRequest identifies the stream whose disk usage is requested.
type GetStreamUsageRequest struct {
	// Stream is the name of the stream to measure.
	Stream string `path:"stream" minLength:"1"`
}

// GetStreamUsageResponse carries the measured storage footprint.
type GetStreamUsageResponse struct {
	// Success reports whether the measurement was returned.
	Success bool `json:"success"`
	// Usage is the storage footprint of the stream, in bytes.
	Usage int64 `json:"usage"`
}

// NewGetStreamUsageUsecase reports the storage footprint of a stream.
//
// Callers must have the read-streams permission.
func NewGetStreamUsageUsecase(deps GetStreamUsageDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req GetStreamUsageRequest,
				resp *GetStreamUsageResponse,
			) error {
				usage, err := deps.LogStorage.StreamUsage(ctx, req.Stream)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get stream usage",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Usage = usage
				return nil
			},
		),
	)

	u.SetName("get_stream_usage")
	u.SetTitle("Get Stream usage")
	u.SetDescription("Get Stream usage")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetStreamUsageUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}/usage",
	)
}
