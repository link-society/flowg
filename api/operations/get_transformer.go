package operations

import (
	"context"
	"fmt"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	storage "go.link-society.com/flowg/internal/storage/interfaces"
	"go.link-society.com/flowg/pkg/models"
)

// GetTransformerDeps lists the dependencies of [NewGetTransformerUsecase].
type GetTransformerDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// NewGetTransformerUsecase returns the source code of a single transformer.
//
// Callers must have the read-transformers permission. Requesting an unknown
// transformer yields a not-found error.
func NewGetTransformerUsecase(deps GetTransformerDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req schemas.GetTransformerRequest,
				resp *schemas.GetTransformerResponse,
			) error {
				script, err := deps.ConfigStorage.ReadTransformer(ctx, req.Transformer)
				if err != nil {
					deps.Logger.ErrorContext(
						ctx,
						"Failed to get transformer",
						slog.String("transformer", req.Transformer),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if script == nil {
					resp.Success = false
					return status.Wrap(
						fmt.Errorf("transformer %q not found", req.Transformer),
						status.NotFound,
					)
				}

				resp.Success = true
				resp.Script = *script

				return nil
			},
		),
	)

	u.SetName("get_transformer")
	u.SetTitle("Get Transformer")
	u.SetDescription("Get Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewGetTransformerUsecase,
		http.MethodGet,
		"/api/v1/transformers/{transformer}",
	)
}
