package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"
)

type GetStreamUsageRequest struct {
	Stream string `path:"stream" minLength:"1"`
}

type GetStreamUsageResponse struct {
	Success bool  `json:"success"`
	Usage   int64 `json:"usage"`
}

func (ctrl *controller) GetStreamUsageUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req GetStreamUsageRequest,
				resp *GetStreamUsageResponse,
			) error {
				usage, err := ctrl.deps.LogStorage.StreamUsage(ctx, req.Stream)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get stream usage",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Usage = usage
				return nil
			},
		),
	)

	u.SetName("get_stream_usage")
	u.SetTitle("Get Stream usage")
	u.SetDescription("Get Stream usage")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
