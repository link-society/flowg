package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type DeleteTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
}

type DeleteTransformerResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) DeleteTransformerUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req DeleteTransformerRequest,
				resp *DeleteTransformerResponse,
			) error {
				err := ctrl.deps.ConfigStorage.DeleteTransformer(ctx, req.Transformer)
				if err != nil {
					ctrl.logger.ErrorContext(
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
