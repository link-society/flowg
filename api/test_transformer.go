package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"
	"link-society.com/flowg/internal/vrl"
)

type TestTransformerRequest struct {
	Transformer string            `path:"transformer" minLength:"1"`
	Record      map[string]string `json:"record"`
}

type TestTransformerResponse struct {
	Success bool              `json:"success"`
	Record  map[string]string `json:"record"`
}

func TestTransformerUsecase(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req TestTransformerRequest,
				resp *TestTransformerResponse,
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

				resp.Record, err = vrl.ProcessRecord(req.Record, script)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to execute transformer script",
						"channel", "api",
						"transformer", req.Transformer,
						"errorr", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("test_transformer")
	u.SetTitle("Test Transformer Script")
	u.SetDescription("Test Transformer Script")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
