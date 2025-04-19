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

type IngestLogsStructRequest struct {
	Pipeline string              `path:"pipeline" minLength:"1"`
	Records  []map[string]string `json:"records"`
}
type IngestLogsStructResponse struct {
	Success        bool `json:"success"`
	ProcessedCount int  `json:"processed_count"`
}

func (ctrl *controller) IngestLogsStructUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestLogsStructRequest,
				resp *IngestLogsStructResponse,
			) error {
				for _, recordData := range req.Records {
					record := models.NewLogRecord(recordData)
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
				}

				resp.Success = true
				resp.ProcessedCount = len(req.Records)

				return nil
			},
		),
	)

	u.SetName("ingest_logs_struct")
	u.SetTitle("Ingest Structured Logs")
	u.SetDescription("Run structured logs through a pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}
