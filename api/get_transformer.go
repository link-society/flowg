package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/pipelines"
)

type GetTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
}

type GetTransformerResponse struct {
	Success bool   `json:"success"`
	Script  string `json:"script"`
}

func GetTransformerUsecase(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req GetTransformerRequest,
				resp *GetTransformerResponse,
			) error {
				script, err := pipelinesManager.GetTransformerScript(req.Transformer)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get transformer script",
						"channel", "api",
						"transformer", req.Transformer,
						"error", err.Error(),
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
	u.SetTitle("Get Transformer Script")
	u.SetDescription("Get Transformer Script")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound)

	return u
}
