package operations

import (
	"context"
	"fmt"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// GetTransformerDeps lists the dependencies of [NewGetTransformerUsecase].
type GetTransformerDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// GetTransformerRequest identifies the transformer to retrieve.
type GetTransformerRequest struct {
	// Transformer is the name of the transformer to read.
	Transformer string `path:"transformer" minLength:"1"`
}

// GetTransformerResponse carries the source of the requested transformer.
type GetTransformerResponse struct {
	// Success reports whether the transformer was found and returned.
	Success bool `json:"success"`
	// Script is the VRL source code of the transformer.
	Script string `json:"script"`
}

// NewGetTransformerUsecase returns the source code of a single transformer.
//
// Callers must have the read-transformers permission. Requesting an unknown
// transformer yields a not-found error.
func NewGetTransformerUsecase(deps GetTransformerDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req GetTransformerRequest,
				resp *GetTransformerResponse,
			) error {
				script, err := deps.ConfigStorage.ReadTransformer(ctx, req.Transformer)
				if err != nil {
					logger.ErrorContext(
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
