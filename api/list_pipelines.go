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

type ListPipelinesRequest struct{}
type ListPipelinesResponse struct {
	Success   bool     `json:"success"`
	Pipelines []string `json:"pipelines"`
}

func ListPipelinesUsecase(
	authStorage *auth.Storage,
	configStorage *config.Storage,
) usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			authStorage,
			models.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req ListPipelinesRequest,
				resp *ListPipelinesResponse,
			) error {
				pipelines, err := configStorage.ListPipelines(ctx)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list pipelines",
						slog.String("channel", "api"),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				resp.Success = true
				resp.Pipelines = pipelines

				return nil
			},
		),
	)

	u.SetName("list_pipelines")
	u.SetTitle("List Pipelines")
	u.SetDescription("List pipelines")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
