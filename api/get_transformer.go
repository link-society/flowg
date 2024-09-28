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

type GetTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
}

type GetTransformerResponse struct {
	Success bool   `json:"success"`
	Script  string `json:"script"`
}

func GetTransformerUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req GetTransformerRequest,
				resp *GetTransformerResponse,
			) error {
				script, err := configStorage.ReadTransformer(ctx, req.Transformer)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get transformer",
						slog.String("channel", "api"),
						slog.String("transformer", req.Transformer),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				resp.Success = true
				resp.Script = script

				return nil
			},
		),
	)

	u.SetName("get_transformer")
	u.SetTitle("Get Transformer")
	u.SetDescription("Get Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
