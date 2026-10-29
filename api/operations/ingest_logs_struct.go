package operations

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"go.link-society.com/flowg/api/auth"
	"go.link-society.com/flowg/api/routing"
	"go.link-society.com/flowg/api/schemas"

	applog "go.link-society.com/flowg/internal/app/logging"
	"go.link-society.com/flowg/internal/engines/pipelines"
	"go.link-society.com/flowg/internal/models"

	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// IngestLogsStructDeps lists the dependencies of [NewIngestLogsStructUsecase].
type IngestLogsStructDeps struct {
	fx.In

	Logger *slog.Logger

	AuthStorage    storage.AuthStorage
	PipelineRunner pipelines.Runner
}

// NewIngestLogsStructUsecase pushes structured log records through a pipeline.
//
// It is the primary ingestion entry point for callers that already hold
// key/value records. Callers must have the send-logs permission. Ingestion
// stops at the first record that fails. The request is marked sensitive so the
// payload stays out of FlowG's own logs.
func NewIngestLogsStructUsecase(deps IngestLogsStructDeps) usecase.Interactor {
	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.Logger,
			deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req *schemas.IngestLogsStructRequest,
				resp *schemas.IngestLogsStructResponse,
			) error {
				applog.MarkSensitive(ctx)

				for _, recordData := range req.Records {
					record := models.NewLogRecord(recordData)
					err := deps.PipelineRunner.Process(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						record,
					)
					if err != nil {
						if errors.Is(err, &pipelines.PipelineNotFoundError{}) {
							resp.Success = false
							return status.Wrap(err, status.NotFound)
						}

						deps.Logger.DebugContext(
							ctx,
							"Failed to process log entry",
							slog.String("pipeline", req.Pipeline),
							slog.String("error", err.Error()),
						)

						resp.Success = false
						return status.Wrap(err, status.Internal)
					}

					deps.Logger.DebugContext(
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

	u.SetExpectedErrors(status.PermissionDenied, status.NotFound, status.Internal)

	return u
}

func init() {
	routing.RegisterOperation(
		NewIngestLogsStructUsecase,
		http.MethodPost,
		"/api/v1/pipelines/{pipeline}/logs/struct",
	)
}
