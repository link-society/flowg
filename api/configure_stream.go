package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
)

type ConfigureStreamRequest struct {
	Stream string                  `path:"stream" minLength:"1"`
	Config logstorage.StreamConfig `json:"config"`
}

type ConfigureStreamResponse struct {
	Success bool `json:"success"`
}

func ConfigureStreamUsecase(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req ConfigureStreamRequest,
				resp *ConfigureStreamResponse,
			) error {
				err := logDb.ConfigureStream(req.Stream, req.Config)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to configure stream",
						"stream", req.Stream,
						"error", err.Error(),
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
