package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ListStreamsRequest struct{}
type ListStreamsResponse struct {
	Success bool                           `json:"success"`
	Streams map[string]models.StreamConfig `json:"streams"`
}

func (ctrl *controller) ListStreamsUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamsRequest,
				resp *ListStreamsResponse,
			) error {
				streams, err := ctrl.deps.LogStorage.ListStreamConfigs(ctx)
				if err != nil {
					ctrl.logger.ErrorContext(
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
