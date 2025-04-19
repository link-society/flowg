package api

import (
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/x-protobuf":
		ior.logRecords, err = otlp.UnmarshalProtobuf(body)

	case "application/json":
		ior.logRecords, err = otlp.UnmarshalJSON(body)

	default:
		err = fmt.Errorf("unsupported content type: %s", contentType)
	}

	return err
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
