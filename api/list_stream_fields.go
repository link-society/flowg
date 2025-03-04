package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type ListStreamFieldsRequest struct {
	Stream string `path:"stream" minLength:"1"`
}
type ListStreamFieldsResponse struct {
	Success bool     `json:"success"`
	Fields  []string `json:"fields"`
}

func (ctrl *controller) ListStreamFieldsUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamFieldsRequest,
				resp *ListStreamFieldsResponse,
			) error {
				fields, err := ctrl.deps.LogStorage.ListStreamFields(ctx, req.Stream)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to list stream fields",
						slog.String("stream", req.Stream),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Fields = fields

				return nil
			},
		),
	)

	u.SetName("list_stream_fields")
	u.SetTitle("List Stream Fields")
	u.SetDescription("List known fields for a stream")
	u.SetTags("streams")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
