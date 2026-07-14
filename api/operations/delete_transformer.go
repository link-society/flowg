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

// DeleteTransformerDeps lists the dependencies of [NewDeleteTransformerUsecase].
type DeleteTransformerDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewDeleteTransformerUsecase removes a transformer.
//
// Callers must have the write-transformers permission. Deleting an absent
// transformer is treated as a success.
func NewDeleteTransformerUsecase(deps DeleteTransformerDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req schemas.DeleteTransformerRequest,
				resp *schemas.DeleteTransformerResponse,
			) error {
				err := deps.ConfigStorage.DeleteTransformer(ctx, req.Transformer)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to delete transformer",
						slog.String("transformer", req.Transformer),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("delete_transformer")
	u.SetTitle("Delete Transformer")
	u.SetDescription("Delete Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewDeleteTransformerUsecase,
		http.MethodDelete,
		"/api/v1/transformers/{transformer}",
	)
}
