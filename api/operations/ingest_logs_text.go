package operations

import (
	"context"
	"log/slog"
	"time"

	"net/http"
	"strings"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	applog "link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"

	authStorage "link-society.com/flowg/internal/storage/auth"
)

// IngestLogsTextDeps lists the dependencies of [NewIngestLogsTextUsecase].
type IngestLogsTextDeps struct {
	fx.In

	AuthStorage    authStorage.Storage
	PipelineRunner pipelines.Runner
}

// IngestLogsTextRequest carries a plain-text body to push through a pipeline.
type IngestLogsTextRequest struct {
	// Pipeline is the name of the pipeline to run the lines through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// TextBody is the raw text payload; each non-empty line becomes one record.
	TextBody string `contentType:"text/plain"`
}

// IngestLogsTextResponse reports how many lines were processed.
type IngestLogsTextResponse struct {
	// Success reports whether every line was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of lines that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}

// NewIngestLogsTextUsecase ingests a plain-text payload, treating each non-empty
// line as a separate log record.
//
// It suits callers that emit raw log lines rather than structured records; each
// line is wrapped into a record's "content" field. Callers must have the
// send-logs permission. Ingestion stops at the first line that fails. The
// request is marked sensitive so the payload stays out of FlowG's own logs.
func NewIngestLogsTextUsecase(deps IngestLogsTextDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestLogsTextRequest,
				resp *IngestLogsTextResponse,
			) error {
				applog.MarkSensitive(ctx)

				lines := strings.Split(strings.ReplaceAll(req.TextBody, "\r\n", "\n"), "\n")
				var messages []string

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

					err := deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						record,
					)
					if err != nil {
						logger.DebugContext(
							ctx,
							"Failed to process log entry",
							slog.String("pipeline", req.Pipeline),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						return status.Wrap(err, status.Internal)
					}

					logger.DebugContext(
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

	u.SetName("ingest_logs_text")
	u.SetTitle("Ingest Textual Logs")
	u.SetDescription("Run textual logs through a pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewIngestLogsTextUsecase,
		http.MethodPost,
		"/api/v1/pipelines/{pipeline}/logs/text",
	)
}
