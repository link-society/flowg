package api

import (
	"context"
	"log/slog"

	"encoding/json"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
)

type SavePipelineRequest struct {
	Pipeline string           `path:"pipeline" minLength:"1"`
	Flow     config.FlowGraph `json:"flow"`
}

type SavePipelineResponse struct {
	Success bool `json:"success"`
}

func SavePipelineUsecase(
	authDb *auth.Database,
	configStorage *config.Storage,
) usecase.Interactor {
	pipelineSys := config.NewPipelineSystem(configStorage)

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

				err = pipelineSys.Write(req.Pipeline, string(flowData))
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to save pipeline",
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
	u.SetTitle("Save Pipeline")
	u.SetDescription("Save pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.InvalidArgument, status.Internal)

	return u
}
