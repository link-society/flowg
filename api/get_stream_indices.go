package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/models"
	apiUtils "link-society.com/flowg/internal/utils/api"
)

type GetStreamIndicesRequest struct {
	Stream string `path:"stream" minLength:"1"`
}

type GetStreamIndicesResponse struct {
	Success bool                `json:"success"`
	Indices map[string][]string `json:"indices"`
}

func (ctrl *controller) GetStreamIndicesUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req GetStreamIndicesRequest,
				resp *GetStreamIndicesResponse,
			) error {
				indices, err := ctrl.deps.LogStorage.Distinct(ctx, req.Stream)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to get stream indices",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Indices = indices
				return nil
			},
		),
	)

	u.SetName("get_stream_indices")
	u.SetTitle("Get Stream indices")
	u.SetDescription("Get distinct values for indexed fields in a Stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
