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
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// SaveTransformerDeps lists the dependencies of [NewSaveTransformerUsecase].
type SaveTransformerDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	ConfigStorage  storage.ConfigStorage
	PipelineRunner pipelines.Runner
}

// SaveTransformerRequest carries the transformer name and its new source.
type SaveTransformerRequest struct {
	// Transformer is the name of the transformer to create or overwrite.
	Transformer string `path:"transformer" minLength:"1"`
	// Script is the VRL source code to store under that name.
	Script string `json:"script" required:"true"`
}

// SaveTransformerResponse reports the outcome of the save.
type SaveTransformerResponse struct {
	// Success reports whether the transformer was persisted.
	Success bool `json:"success"`
}

// NewSaveTransformerUsecase creates or overwrites a transformer.
//
// Callers must have the write-transformers permission. Persisting a transformer
// invalidates cached pipeline builds so that subsequent runs pick up the new
// source.
func NewSaveTransformerUsecase(deps SaveTransformerDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_WRITE_TRANSFORMERS,
			func(
				ctx context.Context,
				req SaveTransformerRequest,
				resp *SaveTransformerResponse,
			) error {
				if err := deps.ConfigStorage.WriteTransformer(ctx, req.Transformer, req.Script); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to save transformer",
						slog.String("transformer", req.Transformer),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				if err := deps.PipelineRunner.InvalidateAllCachedBuilds(ctx); err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to refresh pipeline cache after save",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("save_transformer")
	u.SetTitle("Save Transformer")
	u.SetDescription("Save Transformer")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewSaveTransformerUsecase,
		http.MethodPut,
		"/api/v1/transformers/{transformer}",
	)
}
