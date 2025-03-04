package api

import (
	"context"
	"log/slog"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/pipelines"
)

type IngestLogRequest struct {
	Pipeline string            `path:"pipeline" minLength:"1"`
	Record   map[string]string `json:"record"`
}
type IngestLogResponse struct {
	Success bool `json:"success"`
}

func (ctrl *controller) IngestLogUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestLogRequest,
				resp *IngestLogResponse,
			) error {
				record := models.NewLogRecord(req.Record)
				err := ctrl.deps.PipelineRunner.Run(
					ctx,
					req.Pipeline,
					pipelines.DIRECT_ENTRYPOINT,
					record,
				)
				if err != nil {
					ctrl.logger.ErrorContext(
						ctx,
						"Failed to process log entry",
						slog.String("pipeline", req.Pipeline),
						slog.String("error", err.Error()),
					)

					resp.Success = false
					return status.Wrap(err, status.Internal)
				}

				ctrl.logger.InfoContext(
					ctx,
					"Log entry processed",
					slog.String("pipeline", req.Pipeline),
				)
				resp.Success = true

				return nil
			},
		),
	)

	u.SetName("ingest_log")
	u.SetTitle("Ingest Log")
	u.SetDescription("Run log record through a pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
