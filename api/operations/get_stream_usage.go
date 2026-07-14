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

// GetStreamUsageDeps lists the dependencies of [NewGetStreamUsageUsecase].
type GetStreamUsageDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
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
				req schemas.GetStreamUsageRequest,
				resp *schemas.GetStreamUsageResponse,
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
