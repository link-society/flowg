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

	authStorage "link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

// ListTransformersDeps lists the dependencies of [NewListTransformersUsecase].
type ListTransformersDeps struct {
	fx.In

	AuthStorage   authStorage.Storage
	ConfigStorage config.Storage
}

// ListTransformersRequest is empty: listing transformers takes no parameters.
type ListTransformersRequest struct{}

// ListTransformersResponse carries the names of the available transformers.
type ListTransformersResponse struct {
	// Success reports whether the listing completed.
	Success bool `json:"success"`
	// Transformers holds the name of every configured transformer.
	Transformers []string `json:"transformers"`
}

// NewListTransformersUsecase enumerates the names of all configured transformers.
//
// Callers must have the read-transformers permission.
func NewListTransformersUsecase(deps ListTransformersDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req ListTransformersRequest,
				resp *ListTransformersResponse,
			) error {
				transformers, err := deps.ConfigStorage.ListTransformers(ctx)
				if err != nil {
					logger.ErrorContext(
						ctx,
						"Failed to list transformers",
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Transformers = transformers

				return nil
			},
		),
	)

	u.SetName("list_transformers")
	u.SetTitle("List Transformers")
	u.SetDescription("List Transformers")
	u.SetTags("transformers")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewListTransformersUsecase,
		http.MethodGet,
		"/api/v1/transformers",
	)
}
