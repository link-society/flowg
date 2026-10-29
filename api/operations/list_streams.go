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

// ListStreamsDeps lists the dependencies of [NewListStreamsUsecase].
type ListStreamsDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewListStreamsUsecase enumerates all known streams with their configurations.
//
// Callers must have the read-streams permission.
func NewListStreamsUsecase(deps ListStreamsDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req schemas.ListStreamsRequest,
				resp *schemas.ListStreamsResponse,
			) error {
				streams, err := deps.LogStorage.ListStreamConfigs(ctx)
				if err != nil {
					deps.Logger.ErrorContext(
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
