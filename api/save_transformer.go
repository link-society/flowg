package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/pipelines"
)

type SaveTransformerRequest struct {
	Transformer string `path:"transformer" minLength:"1"`
	Script      string `json:"script"`
}

type SaveTransformerResponse struct {
	Success bool `json:"success"`
}

func SaveTransformerUsecase(pipelinesManager *pipelines.Manager) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req SaveTransformerRequest,
			resp *SaveTransformerResponse,
		) error {
			err := pipelinesManager.SaveTransformerScript(req.Transformer, req.Script)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to save transformer script",
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
	)

	u.SetName("save_transformer")
	u.SetTitle("Save Transformer Script")
	u.SetDescription("Save Transformer Script")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.Internal)

	return u
}
