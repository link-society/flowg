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

// ListStreamsDeps lists the dependencies of [NewListStreamsUsecase].
type ListStreamsDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// ListStreamsRequest is empty: listing streams takes no parameters.
type ListStreamsRequest struct{}

// ListStreamsResponse carries every known stream and its configuration.
type ListStreamsResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Streams maps each stream name to its configuration.
	Streams map[string]models.StreamConfig `json:"streams"`
}

// NewListStreamsUsecase enumerates all known streams with their configurations.
//
// Callers must have the read-streams permission.
func NewListStreamsUsecase(deps ListStreamsDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamsRequest,
				resp *ListStreamsResponse,
			) error {
				streams, err := deps.LogStorage.ListStreamConfigs(ctx)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to list streams",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Streams = streams

				return nil
			},
		),
	)

	u.SetName("list_streams")
	u.SetTitle("List Streams")
	u.SetDescription("List known stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListStreamsUsecase,
		http.MethodGet,
		"/api/v1/streams",
	)
}
