package api

import (
	"context"
	"log/slog"

	"encoding/json"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/auth"
	"link-society.com/flowg/internal/pipelines"
)

type SavePipelineRequest struct {
	Pipeline string              `path:"pipeline" minLength:"1"`
	Flow     pipelines.FlowGraph `json:"flow"`
}

type SavePipelineResponse struct {
	Success bool `json:"success"`
}

func SavePipelineUsecase(
	authDb *auth.Database,
	pipelinesManager *pipelines.Manager,
) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			authDb,
			auth.SCOPE_WRITE_PIPELINES,
			func(
				ctx context.Context,
				req SavePipelineRequest,
				resp *SavePipelineResponse,
			) error {
				flowData, err := json.Marshal(req.Flow)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to marshal pipeline flow",
						"channel", "api",
						"pipeline", req.Pipeline,
						"error", err.Error(),
					)

					resp.Success = false
					return status.Wrap(err, status.InvalidArgument)
				}

				err = pipelinesManager.SavePipelineFlow(req.Pipeline, string(flowData))
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save pipeline flow",
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

	u.SetName("save_pipeline")
	u.SetTitle("Save Pipeline Flow")
	u.SetDescription("Save pipeline flow")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}
