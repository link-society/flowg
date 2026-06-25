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

// ConfigureStreamDeps lists the dependencies of [NewConfigureStreamUsecase].
type ConfigureStreamDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// ConfigureStreamRequest carries the stream name and its new configuration.
type ConfigureStreamRequest struct {
	// Stream is the name of the stream to configure.
	Stream string `path:"stream" minLength:"1"`
	// Config is the retention and indexing configuration to apply.
	Config models.StreamConfig `json:"config" required:"true"`
}

// ConfigureStreamResponse reports the outcome of the configuration change.
type ConfigureStreamResponse struct {
	// Success reports whether the configuration was applied.
	Success bool `json:"success"`
}

// NewConfigureStreamUsecase sets the retention and indexing configuration of a
// stream.
//
// Callers must have the write-streams permission.
func NewConfigureStreamUsecase(deps ConfigureStreamDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req ConfigureStreamRequest,
				resp *ConfigureStreamResponse,
			) error {
				err := deps.LogStorage.ConfigureStream(ctx, req.Stream, req.Config)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to configure stream",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("configure_stream")
	u.SetTitle("Configure Stream")
	u.SetDescription("Configure Stream retention and indexes")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewConfigureStreamUsecase,
		http.MethodPut,
		"/api/v1/streams/{stream}",
	)
}
