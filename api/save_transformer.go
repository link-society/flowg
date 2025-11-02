package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
)

type SaveTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
	Script      string `json:"script" required:"true"`
}

type SaveTransformerResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) SaveTransformerUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req SaveTransformerRequest,
				resp *SaveTransformerResponse,
			) error {
				err := ctrl.deps.ConfigStorage.WriteTransformer(ctx, req.Transformer, req.Script)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to save transformer",
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

	u.SetName("save_transformer")
	u.SetTitle("Save Transformer")
	u.SetDescription("Save Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
