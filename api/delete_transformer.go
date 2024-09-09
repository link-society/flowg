package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type DeleteTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
}

type DeleteTransformerResponse struct {
	Success bool `json:"success"`
}

func DeleteTransformerUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	transformerSys := config.NewTransformerSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req DeleteTransformerRequest,
				resp *DeleteTransformerResponse,
			) error {
				err := transformerSys.Delete(req.Transformer)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to delete transformer",
						"channel", "api",
						"transformer", req.Transformer,
						"error", err.Error(),
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
