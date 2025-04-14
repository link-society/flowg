package otlp

import (
	"errors"
	"log/slog"

	"context"
	"net"
	gohttp "net/http"
	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/utils/proctree"

	collectlogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	collectmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	collecttraces "go.opentelemetry.io/proto/otlp/collector/trace/v1"
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
	rootHandler.Handle("/logs/otlp", h.GetOTLPHandler(ctx, &collectlogs.ExportLogsServiceRequest{}, logsToLogRecords))
	rootHandler.Handle("/metrics/otlp", h.GetOTLPHandler(ctx, &collectmetrics.ExportMetricsServiceRequest{}, metricsToLogRecords))
	rootHandler.Handle("/traces/otlp", h.GetOTLPHandler(ctx, &collecttraces.ExportTraceServiceRequest{}, tracesToLogRecords))

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
