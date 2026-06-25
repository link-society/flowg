package operations

import (
	"context"
	"fmt"
	"log/slog"

	"compress/gzip"
	"io"
	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/rest/request"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"

	applog "link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/otlp"

	"link-society.com/flowg/internal/storage"
)

// IngestLogsOTLPDeps lists the dependencies of [NewIngestLogsOTLPUsecase].
type IngestLogsOTLPDeps struct {
	fx.In

	AuthStorage    storage.AuthStorage
	PipelineRunner pipelines.Runner
}

// IngestLogsOTLPRequest carries an OpenTelemetry logs export to push through a
// pipeline.
//
// It implements [request.Loader] because the OTLP payload may be protobuf or
// JSON and optionally gzip-compressed, which the generic decoder cannot handle
// on its own.
type IngestLogsOTLPRequest struct {
	// Pipeline is the name of the pipeline to run the records through.
	Pipeline string `path:"pipeline" minLength:"1"`
	// ContentEncoding is the payload's transfer encoding; only gzip is accepted.
	ContentEncoding string `header:"Content-Encoding" enum:"gzip" required:"false"`
	// logRecords holds the records decoded from the OTLP payload.
	logRecords []*models.LogRecord

	collectlogs.ExportLogsServiceRequest
}

var _ request.Loader = (*IngestLogsOTLPRequest)(nil)

// LoadFromHTTPRequest decodes the OTLP payload from the raw HTTP request,
// handling optional gzip compression and both protobuf and JSON encodings.
//
// It populates the request's pipeline name and decoded records so the usecase
// can stay agnostic of the wire format.
func (ior *IngestLogsOTLPRequest) LoadFromHTTPRequest(r *http.Request) error {
	ior.Pipeline = r.PathValue("pipeline")
	if ior.Pipeline == "" {
		return fmt.Errorf("pipeline is required")
	}
	defer r.Body.Close()

	slog.InfoContext(
		r.Context(),
		"Parsing OpenTelemetry message",
		slog.String("otlp.content-type", r.Header.Get("Content-Type")),
		slog.String("otlp.content-encoding", r.Header.Get("Content-Encoding")),
	)

	var body []byte

	ior.ContentEncoding = r.Header.Get("Content-Encoding")
	switch ior.ContentEncoding {
	case "gzip":
		// decompress body
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gz.Close()

		data, err := io.ReadAll(gz)
		if err != nil {
			return fmt.Errorf("failed to read gzip body: %w", err)
		}

		body = data

	case "":
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("failed to read raw body: %w", err)
		}

		body = data

	default:
		return fmt.Errorf("unsupported content encoding: %s", ior.ContentEncoding)
	}

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/x-protobuf":
		logRecords, err := otlp.UnmarshalProtobuf(body)
		if err != nil {
			return fmt.Errorf("failed to unmarshal protobuf: %w", err)
		}

		ior.logRecords = logRecords

	case "application/json":
		logRecords, err := otlp.UnmarshalJSON(body)
		if err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}

		ior.logRecords = logRecords

	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	return nil
}

// IngestLogsOTLPResponse reports how many records were processed.
type IngestLogsOTLPResponse struct {
	// Success reports whether every record was processed.
	Success bool `json:"success"`
	// ProcessedCount is the number of records that ran through the pipeline.
	ProcessedCount int `json:"processed_count"`
}

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
				req *IngestLogsOTLPRequest,
				resp *IngestLogsOTLPResponse,
			) error {
				applog.MarkSensitive(ctx)

				for _, logRecord := range req.logRecords {
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
				resp.ProcessedCount = len(req.logRecords)

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
