package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/ffi/vrl"
)

type TestTransformerRequest struct {
	Code   string            `json:"code"`
	Record map[string]string `json:"record"`
}

type TestTransformerResponse struct {
	Success bool              `json:"success"`
	Record  map[string]string `json:"record"`
}

func TestTransformerUsecase(
	authDb *auth.Database,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req TestTransformerRequest,
				resp *TestTransformerResponse,
			) error {
				output, err := vrl.ProcessRecord(req.Record, req.Code)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to execute transformer",
						"channel", "api",
						"errorr", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Record = output

				return nil
			},
		),
	)

	u.SetName("test_transformer")
	u.SetTitle("Test Transformer")
	u.SetDescription("Test Transformer")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
