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

// PurgeStreamDeps lists the dependencies of [NewPurgeStreamUsecase].
type PurgeStreamDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewPurgeStreamUsecase removes a stream along with all of its logs and indexes.
//
// Callers must have the write-streams permission. This is a destructive
// operation: the stored logs cannot be recovered afterwards.
func NewPurgeStreamUsecase(deps PurgeStreamDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req schemas.PurgeStreamRequest,
				resp *schemas.PurgeStreamResponse,
			) error {
				err := deps.LogStorage.DeleteStream(ctx, req.Stream)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to purge stream",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				logger.InfoContext(
					ctx,
					"Log stream purged",
					slog.String("stream", req.Stream),
				)
				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("purge_stream")
	u.SetTitle("Purge Stream")
	u.SetDescription("Remove all logs (and indexes) related to a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewPurgeStreamUsecase,
		http.MethodDelete,
		"/api/v1/streams/{stream}",
	)
}
