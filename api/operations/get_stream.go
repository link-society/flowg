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
)

// GetStreamDeps lists the dependencies of [NewGetStreamUsecase].
type GetStreamDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
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
				req schemas.GetStreamRequest,
				resp *schemas.GetStreamResponse,
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
