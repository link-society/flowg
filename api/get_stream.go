package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
)

type GetStreamRequest struct {
	Stream string `path:"stream" minLength:"1"`
}

type GetStreamResponse struct {
	Success bool                     `json:"success"`
	Config  *logstorage.StreamConfig `json:"config"`
}

func GetStreamUsecase(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) usecase.Interactor {
	metaSys := logstorage.NewMetaSystem(logDb)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetStreamRequest,
				resp *GetStreamResponse,
			) error {
				config, err := metaSys.GetStreamConfig(req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get stream config",
						"channel", "api",
						"stream", req.Stream,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Config = config

				return nil
			},
		),
	)

	u.SetName("get_stream")
	u.SetTitle("Get Stream Configuration")
	u.SetDescription("Get Stream configuration")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
