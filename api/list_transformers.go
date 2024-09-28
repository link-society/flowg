package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
)

type ListTransformersRequest struct{}
type ListTransformersResponse struct {
	Success      bool     `json:"success"`
	Transformers []string `json:"transformers"`
}

func ListTransformersUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_TRANSFORMERS,
			func(
				ctx context.Context,
				req ListTransformersRequest,
				resp *ListTransformersResponse,
			) error {
				transformers, err := configStorage.ListTransformers(ctx)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list transformers",
						slog.String("channel", "api"),
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
