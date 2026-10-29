package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/pkg/models"
)

// ConfigureStreamDeps lists the dependencies of [NewConfigureStreamUsecase].
type ConfigureStreamDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewConfigureStreamUsecase sets the retention and indexing configuration of a
// stream.
//
// Callers must have the write-streams permission.
func NewConfigureStreamUsecase(deps ConfigureStreamDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req schemas.ConfigureStreamRequest,
				resp *schemas.ConfigureStreamResponse,
			) error {
				err := deps.LogStorage.ConfigureStream(ctx, req.Stream, req.Config)
				if err != nil {
					deps.Logger.ErrorContext(
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
