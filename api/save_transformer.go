package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

type SaveTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
	Script      string `json:"script"`
}

type SaveTransformerResponse struct {
	Success bool `json:"success"`
}

func SaveTransformerUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req SaveTransformerRequest,
				resp *SaveTransformerResponse,
			) error {
				err := configStorage.WriteTransformer(ctx, req.Transformer, req.Script)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save transformer",
						slog.String("channel", "api"),
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
