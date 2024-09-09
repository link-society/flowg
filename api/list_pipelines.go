package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type ListPipelinesRequest struct{}
type ListPipelinesResponse struct {
	Success   bool     `json:"success"`
	Pipelines []string `json:"pipelines"`
}

func ListPipelinesUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	pipelineSys := config.NewPipelineSystem(configStorage)

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req ListPipelinesRequest,
				resp *ListPipelinesResponse,
			) error {
				pipelines, err := pipelineSys.List()
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to list pipelines",
						"channel", "api",
						"error", err.Error(),
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
