package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ConfigureStreamRequest struct {
	Stream string              `path:"stream" minLength:"1"`
	Config models.StreamConfig `json:"config" required:"true"`
}

type ConfigureStreamResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) ConfigureStreamUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req ConfigureStreamRequest,
				resp *ConfigureStreamResponse,
			) error {
				err := ctrl.deps.LogStorage.ConfigureStream(ctx, req.Stream, req.Config)
				if err != nil {
					ctrl.logger.ErrorContext(
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
