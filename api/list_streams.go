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

type ListStreamsRequest struct{}
type ListStreamsResponse struct {
	Success bool                           `json:"success"`
	Streams map[string]models.StreamConfig `json:"streams"`
}

func ListStreamsUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamsRequest,
				resp *ListStreamsResponse,
			) error {
				streams, err := logStorage.ListStreamConfigs(ctx)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list streams",
						slog.String("channel", "api"),
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
