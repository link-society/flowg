package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/log"
)

type ConfigureStreamRequest struct {
	Stream string              `path:"stream" minLength:"1"`
	Config models.StreamConfig `json:"config"`
}

type ConfigureStreamResponse struct {
	Success bool `json:"success"`
}

func ConfigureStreamUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req ConfigureStreamRequest,
				resp *ConfigureStreamResponse,
			) error {
				err := logStorage.ConfigureStream(ctx, req.Stream, req.Config)
				if err != nil {
					slog.ErrorContext(
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
