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

type GetStreamRequest struct {
	Stream string `path:"stream" minLength:"1"`
}

type GetStreamResponse struct {
	Success bool                `json:"success"`
	Config  models.StreamConfig `json:"config"`
}

func GetStreamUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req GetStreamRequest,
				resp *GetStreamResponse,
			) error {
				config, err := logStorage.GetOrCreateStreamConfig(ctx, req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get stream config",
						slog.String("channel", "api"),
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Config = config

				return nil
			},
		),
	)

	u.SetName("get_stream")
	u.SetTitle("Get Stream Configuration")
	u.SetDescription("Get or Create Stream configuration")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
