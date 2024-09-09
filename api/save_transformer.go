package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type SaveTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
	Script      string `json:"script"`
}

type SaveTransformerResponse struct {
	Success bool `json:"success"`
}

func SaveTransformerUsecase(
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
				req SaveTransformerRequest,
				resp *SaveTransformerResponse,
			) error {
				err := transformerSys.Write(req.Transformer, req.Script)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save transformer",
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

	u.SetName("save_transformer")
	u.SetTitle("Save Transformer")
	u.SetDescription("Save Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
