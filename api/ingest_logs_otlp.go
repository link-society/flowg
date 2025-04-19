package api

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/swaggest/rest/request"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	apiUtils "link-society.com/flowg/internal/utils/api"
	"link-society.com/flowg/internal/utils/otlp"

	"link-society.com/flowg/internal/models"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"link-society.com/flowg/internal/engines/pipelines"
)

type IngestLogsOTLPRequest struct {
	Pipeline   string `path:"pipeline" minLength:"1"`
	logRecords []*models.LogRecord
	collectlogs.ExportLogsServiceRequest
}

var _ request.Loader = (*IngestLogsOTLPRequest)(nil)

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

	contentEncoding := r.Header.Get("Content-Encoding")
	switch contentEncoding {
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
		return fmt.Errorf("unsupported content encoding: %s", contentEncoding)
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

type IngestLogsOTLPResponse struct {
	Success        bool `json:"success"`
	ProcessedCount int  `json:"processed_count"`
}

func (ctrl *controller) IngestLogsOTLPUsecase() usecase.Interactor {
	u := usecase.NewInteractor(
		apiUtils.RequireScopeApiDecorator(
			ctrl.deps.AuthStorage,
			models.SCOPE_SEND_LOGS,
			func(
				ctx context.Context,
				req *IngestLogsOTLPRequest,
				resp *IngestLogsOTLPResponse,
			) error {
				for _, logRecord := range req.logRecords {
					err := ctrl.deps.PipelineRunner.Run(
						ctx,
						req.Pipeline,
						pipelines.DIRECT_ENTRYPOINT,
						logRecord,
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
