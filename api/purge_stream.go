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

type PurgeStreamRequest struct {
	Stream string `path:"stream" minLength:"1"`
}
type PurgeStreamResponse struct {
	Success bool `json:"success"`
}

func PurgeStreamUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_STREAMS,
			func(
				ctx context.Context,
				req PurgeStreamRequest,
				resp *PurgeStreamResponse,
			) error {
				err := logStorage.DeleteStream(ctx, req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to purge stream",
						slog.String("channel", "api"),
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				slog.InfoContext(
					ctx,
					"Log stream purged",
					slog.String("channel", "api"),
					slog.String("stream", req.Stream),
				)
				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("purge_stream")
	u.SetTitle("Purge Stream")
	u.SetDescription("Remove all logs (and indexes) related to a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
