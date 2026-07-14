package operations

import (
	"context"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/api/schemas"

	applog "link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"

	storage "link-society.com/flowg/internal/storage/interfaces"
)

// IngestLogsOTLPDeps lists the dependencies of [NewIngestLogsOTLPUsecase].
type IngestLogsOTLPDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	PipelineRunner pipelines.Runner
}

var _ request.Loader = (*schemas.IngestLogsOTLPRequest)(nil)

// NewIngestLogsOTLPUsecase ingests an OpenTelemetry logs export, pushing each
// decoded record through a pipeline.
//
// It lets OpenTelemetry-instrumented systems ship logs to FlowG natively.
// Callers must have the send-logs permission. Ingestion stops at the first
// record that fails. The request is marked sensitive so the payload stays out
// of FlowG's own logs.
func NewIngestLogsOTLPUsecase(deps IngestLogsOTLPDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopeApiDecorator(
			deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req *schemas.IngestLogsOTLPRequest,
				resp *schemas.IngestLogsOTLPResponse,
			) error {
				applog.MarkSensitive(ctx)

				for _, logRecord := range req.LogRecords {
					err := deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						logRecord,
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
				resp.ProcessedCount = len(req.LogRecords)

				return nil
			},
		),
	)

	u.SetName("ingest_logs_otlp")
	u.SetTitle("Ingest OpenTelemetry logs")

	u.SetDescription("Run OpenTelemetry logs through a pipeline")
	u.SetTags("pipelines")

	u.SetExpectedErrors(status.PermissionDenied, status.Internal)

	return u
}

// annotateIngestLogsOTLP documents the request body as a binary OTLP protobuf payload.
func annotateIngestLogsOTLP(oc openapi.OperationContext) error {
	oc.AddReqStructure(nil, func(cu *openapi.ContentUnit) {
		cu.ContentType = "application/x-protobuf"
		cu.Format = "binary"
		cu.Description = "OpenTelemetry Export Logs Service Request"
	})

	return nil
}

func init() {
	routing.RegisterOperation(
		NewIngestLogsOTLPUsecase,
		http.MethodPost,
		"/api/v1/pipelines/{pipeline}/logs/otlp",
		routing.Annotated(annotateIngestLogsOTLP),
	)
}
