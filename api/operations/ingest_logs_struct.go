package operations

import (
	"context"
	"log/slog"

	"net/http"

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

// IngestLogsStructDeps lists the dependencies of [NewIngestLogsStructUsecase].
type IngestLogsStructDeps struct {
	fx.In

	AuthStorage    authStorage.Storage
	PipelineRunner pipelines.Runner
}

// IngestLogsStructRequest carries structured records to push through a pipeline.
type IngestLogsStructRequest struct {
	// Pipeline is the name of the pipeline to run the records through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// Records are the structured log records to ingest.
	Records []map[string]string `json:"records" required:"true"`
}

// IngestLogsStructResponse reports how many records were processed.
type IngestLogsStructResponse struct {
	// Success reports whether every record was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of records that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}

// NewIngestLogsStructUsecase pushes structured log records through a pipeline.
//
// It is the primary ingestion entry point for callers that already hold
// key/value records. Callers must have the send-logs permission. Ingestion
// stops at the first record that fails. The request is marked sensitive so the
// payload stays out of FlowG's own logs.
func NewIngestLogsStructUsecase(deps IngestLogsStructDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req IngestLogsStructRequest,
				resp *IngestLogsStructResponse,
			) error {
				applog.MarkSensitive(ctx)

				for _, recordData := range req.Records {
					record := models.NewLogRecord(recordData)
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

func init() {
	routing.RegisterOperation(
		NewIngestLogsStructUsecase,
		http.MethodPost,
		"/api/v1/pipelines/{pipeline}/logs/struct",
	)
}
