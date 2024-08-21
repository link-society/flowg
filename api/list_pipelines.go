package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/pipelines"
)

type ListPipelinesRequest struct{}
type ListPipelinesResponse struct {
	Success   bool     `json:"success"`
	Pipelines []string `json:"pipelines"`
}

func ListPipelinesUsecase(pipelinesManager *pipelines.Manager) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req ListPipelinesRequest,
			resp *ListPipelinesResponse,
		) error {
			pipelines, err := pipelinesManager.ListPipelines()
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
	)

	u.SetName("list_pipelines")
	u.SetTitle("List Pipelines")
	u.SetDescription("List pipelines")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.Internal)

	return u
}
