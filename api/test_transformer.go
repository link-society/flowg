package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/ffi/vrl"
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
	configStorage *config.Storage,
) usecase.Interactor {
	transformerSys := config.NewTransformerSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req TestTransformerRequest,
				resp *TestTransformerResponse,
			) error {
				script, err := transformerSys.Read(req.Transformer)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get transformer",
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
						"Failed to execute transformer",
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
	u.SetTitle("Test Transformer")
	u.SetDescription("Test Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
