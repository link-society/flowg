package otlp

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"sync"

	"context"
	"net"
	"net/http"
	gohttp "net/http"
	"time"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	collectmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collecttraces "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"google.golang.org/protobuf/proto"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	logger *slog.Logger

	opts   *ServerOptions
	server *gohttp.Server

	setupShutdown func(context.Context) error
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	var err error

	h.setupShutdown, err = setupOTelSDK(ctx)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to setup OpenTelemetry SDK",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	rootHandler := gohttp.NewServeMux()
	rootHandler.Handle("/logs/otlp", h)
	rootHandler.Handle("/metrics/otlp", h)
	rootHandler.Handle("/traces/otlp", h)

	h.server = &gohttp.Server{
		Addr:      h.opts.BindAddress,
		Handler:   logging.NewMiddleware(rootHandler),
		TLSConfig: h.opts.TlsConfig,
	}

	h.logger.InfoContext(ctx, "Starting HTTP server")

	l, err := net.Listen("tcp", h.opts.BindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to start HTTP server",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	if h.opts.TlsConfig != nil {
		go h.server.ServeTLS(l, "", "")
	} else {
		go h.server.Serve(l)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	h.logger.InfoContext(ctx, "Stopping OTLP server")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	stopErr := h.server.Shutdown(ctx)

	if stopErr != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to shutdown OTLP server",
			slog.String("error", err.Error()),
		)

		err = errors.Join(err, stopErr)
	}

	stopErr = h.setupShutdown(ctx)
	if stopErr != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to shutdown OTLP SDK",
			slog.String("error", stopErr.Error()),
		)
		err = errors.Join(err, stopErr)
	}

	return err
}

func (h *procHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	logRecords := make([]models.LogRecord, 0)

	switch r.URL.Path {
	case "/logs/otlp":
		req := &collectlogs.ExportLogsServiceRequest{}
		if err := proto.Unmarshal(body, req); err != nil {
			http.Error(w, "invalid protobuf", 400)
			return
		}

		for i, resourceLogs := range req.GetResourceLogs() {
			log.Printf("ResourceLogs %d: %v", i, resourceLogs)
			for j, scopeLogs := range resourceLogs.GetScopeLogs() {
				log.Printf("ScopeLogs %d: %v", j, scopeLogs)
				for k, logRecord := range scopeLogs.GetLogRecords() {
					log.Printf("LogRecord %d: %v", k, logRecord)

					logRecordModel, err := LogToLogRecord(logRecord)
					if err != nil {
						log.Printf("Error converting logRecord to LogRecord: %v", err)
						continue
					}
					logRecords = append(logRecords, logRecordModel)
				}
			}
		}
	case "/traces/otlp":
		req := &collecttraces.ExportTraceServiceRequest{}
		if err := proto.Unmarshal(body, req); err != nil {
			http.Error(w, "invalid protobuf", 400)
			return
		}
		logRecords := make([]models.LogRecord, 0)
		for i, resourceSpan := range req.ResourceSpans {
			log.Printf("ResourceLogs %d: %v", i, resourceSpan)
			for j, scopeSpan := range resourceSpan.GetScopeSpans() {
				log.Printf("ScopeLogs %d: %v", j, scopeSpan)
				for k, span := range scopeSpan.GetSpans() {
					log.Printf("LogRecord %d: %v", k, span)

					logRecordModel, err := SpanToLogRecord(span)
					if err != nil {
						log.Printf("Error converting logRecord to LogRecord: %v", err)
						continue
					}
					logRecords = append(logRecords, logRecordModel)
				}
			}
		}
	case "/metrics/otlp":
		req := &collectmetrics.ExportMetricsServiceRequest{}
		if err := proto.Unmarshal(body, req); err != nil {
			http.Error(w, "invalid protobuf", 400)
			return
		}
		logRecords := make([]models.LogRecord, 0)
		for i, resourceMetrics := range req.GetResourceMetrics() {
			log.Printf("ResourceMetrics %d: %v", i, resourceMetrics)
			for j, scopeMetrics := range resourceMetrics.GetScopeMetrics() {
				log.Printf("ScopeMetrics %d: %v", j, scopeMetrics)
				for k, metric := range scopeMetrics.GetMetrics() {
					log.Printf("Metric %d: %v", k, metric)

					logRecordModel, err := MetricToLogRecord(metric)
					if err != nil {
						log.Printf("Error converting metric to LogRecord: %v", err)
						continue
					}
					logRecords = append(logRecords, logRecordModel)
				}
			}
		}
	default:
		http.Error(w, "Unsupported path", http.StatusNotFound)
		return
	}

	err = h.sendToPipelines(r.Context(), logRecords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
