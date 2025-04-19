package api

import (
	"context"
	"log/slog"

	"strings"
	"time"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"

	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/engines/pipelines"
)

type IngestLogsTextRequest struct {
	Pipeline string `path:"pipeline" minLength:"1"`
	TextBody string `contentType:"text/plain"`
}
type IngestLogsTextResponse struct {
	Success        bool `json:"success"`
	ProcessedCount int  `json:"processed_count"`
}

func (ctrl *controller) IngestLogsTextUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestLogsTextRequest,
				resp *IngestLogsTextResponse,
			) error {
				lines := strings.Split(strings.ReplaceAll(req.TextBody, "\r\n", "\n"), "\n")
				messages := []string{}

				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" {
						messages = append(messages, line)
					}
				}

				for _, message := range messages {
					record := &models.LogRecord{
						Timestamp: time.Now(),
						Fields: map[string]string{
							"content": message,
						},
					}

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
				resp.ProcessedCount = len(messages)

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
