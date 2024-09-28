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

type ListStreamFieldsRequest struct {
	Stream string `path:"stream" minLength:"1"`
}
type ListStreamFieldsResponse struct {
	Success bool     `json:"success"`
	Fields  []string `json:"fields"`
}

func ListStreamFieldsUsecase(
	authStorage *auth.Storage,
	logStorage *log.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamFieldsRequest,
				resp *ListStreamFieldsResponse,
			) error {
				fields, err := logStorage.ListStreamFields(ctx, req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list stream fields",
						slog.String("channel", "api"),
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
