package otlp

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"context"
	"net/http"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	collectmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collecttraces "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
)

type requestToLogRecords = func(request proto.Message) ([]models.LogRecord, error)

func logsToLogRecords(request proto.Message) ([]models.LogRecord, error) {
	logRecords := make([]models.LogRecord, 0)

	req := request.(*collectlogs.ExportLogsServiceRequest)

	for _, resourceLogs := range req.GetResourceLogs() {
		for _, scopeLogs := range resourceLogs.GetScopeLogs() {
			for _, logRecord := range scopeLogs.GetLogRecords() {

				logRecordModel, err := LogToLogRecord(logRecord)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}

	return logRecords, nil
}

func tracesToLogRecords(request proto.Message) ([]models.LogRecord, error) {
	logRecords := make([]models.LogRecord, 0)

	req := request.(*collecttraces.ExportTraceServiceRequest)

	for _, resourceSpan := range req.ResourceSpans {
		for _, scopeSpan := range resourceSpan.GetScopeSpans() {
			for _, span := range scopeSpan.GetSpans() {
				logRecordModel, err := SpanToLogRecord(span)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}
	return logRecords, nil
}
func metricsToLogRecords(request proto.Message) ([]models.LogRecord, error) {
	logRecords := make([]models.LogRecord, 0)

	req := request.(*collectmetrics.ExportMetricsServiceRequest)

	for _, resourceMetrics := range req.GetResourceMetrics() {
		for _, scopeMetrics := range resourceMetrics.GetScopeMetrics() {
			for _, metric := range scopeMetrics.GetMetrics() {
				logRecordModel, err := MetricToLogRecord(metric)
				if err != nil {
					return nil, err
				}
				logRecords = append(logRecords, logRecordModel)
			}
		}
	}
	return logRecords, nil
}

func (h *procHandler) sendToPipelines(ctx context.Context, logRecords []models.LogRecord) error {
	pipelineNames, err := h.opts.ConfigStorage.ListPipelines(ctx)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to list pipelines",
			slog.String("error", err.Error()),
		)
		return err
	}

	wg := sync.WaitGroup{}

	for _, pipelineName := range pipelineNames {
		wg.Add(1)
		go func(pipelineName string) {
			defer wg.Done()

			for _, logRecord := range logRecords {

				err := h.opts.PipelineRunner.Run(
					ctx,
					pipelineName,
					pipelines.SYSLOG_ENTRYPOINT,
					&logRecord,
				)
				if err != nil {
					h.logger.ErrorContext(
						ctx,
						"Failed to process log entry",
						slog.String("pipeline", pipelineName),
						slog.String("error", err.Error()),
					)
				}
			}
		}(pipelineName)
	}

	wg.Wait()

	return nil
}

func (h *procHandler) GetOTLPHandler(ctx context.Context, request proto.Message, logRecordsGetter requestToLogRecords) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "HTTP METHOD POST only", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/x-protobuf":
			err := proto.Unmarshal(body, request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case "application/json":
			err := json.Unmarshal(body, request)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, fmt.Sprintf("Unsupported content type: %s", contentType), http.StatusUnsupportedMediaType)
			return
		}

		logRecords, err := logRecordsGetter(request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.sendToPipelines(r.Context(), logRecords)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
