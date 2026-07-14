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
	"link-society.com/flowg/api/schemas"

	"link-society.com/flowg/internal/models"
	storage "link-society.com/flowg/internal/storage/interfaces"
	"link-society.com/flowg/internal/utils/langs/vrl"
)

// TestTransformerDeps lists the dependencies of [NewTestTransformerUsecase].
type TestTransformerDeps struct {
	fx.In

	AuthStorage storage.AuthStorage
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
				req schemas.TestTransformerRequest,
				resp *schemas.TestTransformerResponse,
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
