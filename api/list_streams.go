package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
)

type ListStreamsRequest struct{}
type ListStreamsResponse struct {
	Success bool     `json:"success"`
	Streams []string `json:"streams"`
}

func ListStreamsUsecase(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamsRequest,
				resp *ListStreamsResponse,
			) error {
				streams, err := logDb.ListStreams()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list streams",
						"channel", "api",
						"error", err.Error(),
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
