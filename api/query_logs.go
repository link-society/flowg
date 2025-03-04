package api

import (
	"context"
	"log/slog"

	"time"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/ffi/filterdsl"
)

type QueryStreamRequest struct {
	Stream string    `path:"stream" minLength:"1"`
	From   time.Time `query:"from" format:"date-time" required:"true"`
	To     time.Time `query:"to" format:"date-time" required:"true"`
	Filter *string   `query:"filter"`
}

type QueryStreamResponse struct {
	Success bool               `json:"success"`
	Records []models.LogRecord `json:"records"`
}

func (ctrl *controller) QueryStreamUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req QueryStreamRequest,
				resp *QueryStreamResponse,
			) error {
				var filter filterdsl.Filter

				if req.Filter != nil {
					var err error
					filter, err = filterdsl.Compile(*req.Filter)
					if err != nil {
						ctrl.logger.ErrorContext(
							ctx,
							"Failed to compile filter",
							slog.String("stream", req.Stream),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						resp.Records = nil
						return status.Wrap(err, status.InvalidArgument)
					}
				} else {
					filter = nil
				}

				records, err := ctrl.deps.LogStorage.FetchLogs(ctx, req.Stream, req.From, req.To, filter)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to query logs",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Records = records
				return nil
			},
		),
	)

	u.SetName("query_stream")
	u.SetTitle("Query Stream")
	u.SetDescription("Query logs from a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}
