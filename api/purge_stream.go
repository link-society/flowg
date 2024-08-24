package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/logstorage"
)

type PurgeStreamRequest struct {
	Stream string `path:"stream" minLength:"1"`
}
type PurgeStreamResponse struct {
	Success bool `json:"success"`
}

func PurgeStreamUsecase(db *logstorage.Storage) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req PurgeStreamRequest,
			resp *PurgeStreamResponse,
		) error {
			err := db.Purge(ctx, req.Stream)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to purge stream",
					"channel", "api",
					"stream", req.Stream,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			slog.InfoContext(
				ctx,
				"Log stream purged",
				"channel", "api",
				"stream", req.Stream,
			)
			resp.Success = true

			return nil
		},
	)

	u.SetName("purge_stream")
	u.SetTitle("Purge Stream")
	u.SetDescription("Remove all logs (and indexes) related to a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.Internal)

	return u
}
