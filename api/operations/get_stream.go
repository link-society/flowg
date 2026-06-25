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

	"link-society.com/flowg/internal/storage"
)

// GetStreamDeps lists the dependencies of [NewGetStreamUsecase].
type GetStreamDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// GetStreamRequest identifies the stream whose configuration is requested.
type GetStreamRequest struct {
	// Stream is the name of the stream to read.
	Stream string `path:"stream" minLength:"1"`
}

// GetStreamResponse carries the configuration of the requested stream.
type GetStreamResponse struct {
	// Success reports whether the configuration was returned.
	Success bool `json:"success"`
	// Config is the stream's retention and indexing configuration.
	Config models.StreamConfig `json:"config"`
}

// NewGetStreamUsecase returns the configuration of a stream, creating it with
// defaults if it does not yet exist.
//
// Callers must have the read-streams permission.
func NewGetStreamUsecase(deps GetStreamDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req GetStreamRequest,
				resp *GetStreamResponse,
			) error {
				config, err := deps.LogStorage.GetOrCreateStreamConfig(ctx, req.Stream)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get stream config",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Config = config

				return nil
			},
		),
	)

	u.SetName("get_stream")
	u.SetTitle("Get Stream Configuration")
	u.SetDescription("Get or Create Stream configuration")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetStreamUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}",
	)
}
