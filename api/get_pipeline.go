package api

import (
	"context"
	"log/slog"

	"encoding/json"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/pipelines"
)

type GetPipelineRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
}

type GetPipelineResponse struct {
	Success bool                `json:"success"`
	Flow    pipelines.FlowGraph `json:"flow"`
}

func GetPipelineUsecase(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_READ_PIPELINES,
			func(
				ctx context.Context,
				req GetPipelineRequest,
				resp *GetPipelineResponse,
			) error {
				flowData, err := pipelinesManager.GetPipelineFlow(req.Pipeline)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to get pipeline flow",
						"channel", "api",
						"pipeline", req.Pipeline,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.NotFound)
				}

				if err := json.Unmarshal([]byte(flowData), &resp.Flow); err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to decode pipeline flow",
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
		),
	)

	u.SetName("get_pipeline")
	u.SetTitle("Get Pipeline Flow")
	u.SetDescription("Get pipeline flow")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}
