package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/logstorage"
)

type ListStreamFieldsRequest struct {
	Stream string `path:"stream" minLength:"1"`
}
type ListStreamFieldsResponse struct {
	Success bool     `json:"success"`
	Fields  []string `json:"fields"`
}

func ListStreamFieldsUsecase(
	authDb *auth.Database,
	logDb *logstorage.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req ListStreamFieldsRequest,
				resp *ListStreamFieldsResponse,
			) error {
				fields, err := logDb.ListStreamFields(req.Stream)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list stream fields",
						"channel", "api",
						"stream", req.Stream,
						"error", err.Error(),
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
