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

// GetStreamIndicesDeps lists the dependencies of [NewGetStreamIndicesUsecase].
type GetStreamIndicesDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewGetStreamIndicesUsecase returns the distinct values held by each indexed
// field of a stream.
//
// It backs query builders that offer the known values of an index. Callers must
// have the read-streams permission.
func NewGetStreamIndicesUsecase(deps GetStreamIndicesDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req schemas.GetStreamIndicesRequest,
				resp *schemas.GetStreamIndicesResponse,
			) error {
				indices, err := deps.LogStorage.Distinct(ctx, req.Stream)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to get stream indices",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Indices = indices
				return nil
			},
		),
	)

	u.SetName("get_stream_indices")
	u.SetTitle("Get Stream indices")
	u.SetDescription("Get distinct values for indexed fields in a Stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetStreamIndicesUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}/indices",
	)
}
