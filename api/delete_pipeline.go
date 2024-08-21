package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/pipelines"
)

type DeletePipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type DeletePipelineResponse struct {
	Success bool `json:"success"`
}

func DeletePipelineUsecase(pipelinesManager *pipelines.Manager) usecase.Interactor {
	u := usecase.NewInteractor(
		func(
			ctx context.Context,
			req DeletePipelineRequest,
			resp *DeletePipelineResponse,
		) error {
			err := pipelinesManager.DeletePipelineFlow(req.Pipeline)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"Failed to delete pipeline flow",
					"channel", "api",
					"pipeline", req.Pipeline,
					"error", err.Error(),
				)

				resp.Success = false
				return status.Wrap(err, status.Internal)
			}

			resp.Success = true

			return nil
		},
	)

	u.SetName("delete_pipeline")
	u.SetTitle("Delete Pipeline Flow")
	u.SetDescription("Delete pipeline flow")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.Internal)

	return u
}
