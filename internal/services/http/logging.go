package http

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"bytes"
	gohttp "net/http"

	applog "link-society.com/flowg/internal/app/logging"
)

// loggingMiddleware is the access-log middleware wrapping the root handler.
type loggingMiddleware struct {
	handler gohttp.Handler
}

var _ gohttp.Handler = (*loggingMiddleware)(nil)

// loggingResponseWriter wraps the real ResponseWriter to capture the status code
// and, in verbose mode, buffer the response body for dumping on failure.
type loggingResponseWriter struct {
	ctx        context.Context
	parent     gohttp.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

var _ gohttp.ResponseWriter = (*loggingResponseWriter)(nil)
var _ gohttp.Flusher = (*loggingResponseWriter)(nil)

// newLoggingMiddleware wraps handler with the HTTP access-log middleware.
//
// For each request it assigns a correlation id (taken from the
// "X-Correlation-Id" header, or generated when absent), propagates it through
// the request context, and emits a structured access log record on the
// "accesslog" channel once the wrapped handler returns. Requests flagged as
// sensitive are logged at debug level instead of info.
//
// When applog.VERBOSE_LOGGING is enabled, response bodies of failed requests
// (status >= 400) are buffered and dumped to standard error.
func newLoggingMiddleware(handler gohttp.Handler) gohttp.Handler {
	return &loggingMiddleware{handler: handler}
}

// ServeHTTP assigns a correlation id, runs the wrapped handler and emits the
// access-log record (at debug level for requests marked sensitive).
func (m *loggingMiddleware) ServeHTTP(w gohttp.ResponseWriter, r *gohttp.Request) {
	correlationId := r.Header.Get("X-Correlation-Id")
	ctx := applog.WithSensitiveMarker(applog.WithCorrelationId(r.Context(), correlationId))
	req := r.WithContext(ctx)
	resp := &loggingResponseWriter{
		ctx:        ctx,
		parent:     w,
		buf:        bytes.NewBuffer(nil),
		statusCode: gohttp.StatusOK,
	}
	m.handler.ServeHTTP(resp, req)

	logFn := slog.InfoContext
	if applog.IsMarkedSensitive(req.Context()) {
		logFn = slog.DebugContext
	}

	logFn(
		req.Context(),
		"http request",
		slog.String("channel", "accesslog"),
		slog.Group("http",
			slog.Group("req",
				slog.String("method", req.Method),
				slog.String("url", req.URL.Path),
			),
			slog.Group("resp",
				slog.Int("status", resp.statusCode),
			),
		),
	)
}

func (w *loggingResponseWriter) Header() gohttp.Header {
	return w.parent.Header()
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if applog.VERBOSE_LOGGING {
		w.buf.Write(b)
	}

	return w.parent.Write(b)
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.parent.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Flush() {
	if applog.VERBOSE_LOGGING && w.statusCode >= 400 {
		correlationId := w.ctx.Value(applog.CORRELATION_ID).(string)
		fmt.Fprintf(os.Stderr, "---begin: %s---\n", correlationId)
		fmt.Fprintln(os.Stderr, w.buf.String())
		fmt.Fprintf(os.Stderr, "---end: %s---\n", correlationId)
		w.buf.Reset()
	}

	if f, ok := w.parent.(gohttp.Flusher); ok {
		f.Flush()
	}
}
