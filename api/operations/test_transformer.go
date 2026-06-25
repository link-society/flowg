package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/langs/vrl"

	"link-society.com/flowg/internal/storage"
)

// TestTransformerDeps lists the dependencies of [NewTestTransformerUsecase].
type TestTransformerDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
}

// TestTransformerRequest carries the transformer source and a sample record to
// run it against.
type TestTransformerRequest struct {
	// Code is the VRL source to evaluate, without persisting it.
	Code string `json:"code" required:"true"`
	// Record is the input log record fed to the transformer.
	Record map[string]string `json:"record" required:"true"`
}

// TestTransformerResponse carries the records produced by the trial run.
type TestTransformerResponse struct {
	// Success reports whether the script compiled and ran.
	Success bool `json:"success"`
	// Records holds the output records emitted by the transformer.
	Records []map[string]string `json:"records"`
}

// NewTestTransformerUsecase evaluates a transformer against a sample record
// without storing it.
//
// It lets authors validate VRL source before saving it. Callers must have the
// read-transformers permission. Source that fails to compile or run is reported as
// an unprocessable-entity error rather than a server fault.
func NewTestTransformerUsecase(deps TestTransformerDeps) usecase.Interactor {
	logger := logging.Logger()

	const UnprocessableEntityCode = 422

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req TestTransformerRequest,
				resp *TestTransformerResponse,
			) error {
				runner, err := vrl.NewScriptRunner(req.Code)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to compile transformer",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return usecase.Error{
						AppCode:    UnprocessableEntityCode,
						StatusCode: status.InvalidArgument,
						Value:      err,
					}
				}
				defer runner.Close()

				output, err := runner.TransformLog(req.Record)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to execute transformer",
						slog.String("error", err.Error()),
					)

					resp.Success = false

					return usecase.Error{
						AppCode:    UnprocessableEntityCode,
						StatusCode: status.InvalidArgument,
						Value:      err,
					}
				}

				resp.Success = true
				resp.Records = output

				return nil
			},
		),
	)

	u.SetName("test_transformer")
	u.SetTitle("Test Transformer")
	u.SetDescription("Test Transformer")
	u.SetTags("tests")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal, status.FailedPrecondition)

	return u
}

func init() {
	routing.RegisterOperation(
		NewTestTransformerUsecase,
		http.MethodPost,
		"/api/v1/test/transformer",
	)
}
