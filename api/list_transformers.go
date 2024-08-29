package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/pipelines"
)

type ListTransformersRequest struct{}
type ListTransformersResponse struct {
	Success      bool     `json:"success"`
	Transformers []string `json:"transformers"`
}

func ListTransformersUsecase(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req ListTransformersRequest,
				resp *ListTransformersResponse,
			) error {
				transformers, err := pipelinesManager.ListTransformers()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list transformers",
						"channel", "api",
						"error", err.Error(),
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
