package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// ListStreamFieldsDeps lists the dependencies of [NewListStreamFieldsUsecase].
type ListStreamFieldsDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
	LogStorage  storage.LogStorage
}

// NewListStreamFieldsUsecase enumerates the field names observed across a stream's
// records.
//
// It backs query builders that suggest available fields. Callers must have the
// read-streams permission.
func NewListStreamFieldsUsecase(deps ListStreamFieldsDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_STREAMS,
			func(
				ctx context.Context,
				req schemas.ListStreamFieldsRequest,
				resp *schemas.ListStreamFieldsResponse,
			) error {
				fields, err := deps.LogStorage.ListStreamFields(ctx, req.Stream)
				if err != nil {
					logger.ErrorContext(
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

func init() {
	routing.RegisterOperation(
		NewListStreamFieldsUsecase,
		http.MethodGet,
		"/api/v1/streams/{stream}/fields",
	)
}
