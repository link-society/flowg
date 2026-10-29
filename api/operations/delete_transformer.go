package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/pkg/models"
)

// DeleteTransformerDeps lists the dependencies of [NewDeleteTransformerUsecase].
type DeleteTransformerDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewDeleteTransformerUsecase removes a transformer.
//
// Callers must have the write-transformers permission. Deleting an absent
// transformer is treated as a success.
func NewDeleteTransformerUsecase(deps DeleteTransformerDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req schemas.DeleteTransformerRequest,
				resp *schemas.DeleteTransformerResponse,
			) error {
				err := deps.ConfigStorage.DeleteTransformer(ctx, req.Transformer)
				if err != nil {
					deps.Logger.ErrorContext(
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
