package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"bytes"
	"net/http"
)

type middleware struct {
	handler http.Handler
}

var _ http.Handler = (*middleware)(nil)

type responseWriter struct {
	ctx        context.Context
	parent     http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

var _ http.ResponseWriter = (*responseWriter)(nil)
var _ http.Flusher = (*responseWriter)(nil)

func NewMiddleware(handler http.Handler) http.Handler {
	return &middleware{handler: handler}
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	correlationId := r.Header.Get("X-Correlation-Id")
	ctx := WithSensitiveMarker(WithCorrelationId(r.Context(), correlationId))
	req := r.WithContext(ctx)
	resp := &responseWriter{
		ctx:        ctx,
		parent:     w,
		buf:        bytes.NewBuffer(nil),
		statusCode: http.StatusOK,
	}
	m.handler.ServeHTTP(resp, req)

	logFn := slog.InfoContext
	if IsMarkedSensitive(req.Context()) {
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

func (w *responseWriter) Header() http.Header {
	return w.parent.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if VERBOSE_LOGGING {
		w.buf.Write(b)
	}

	return w.parent.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.parent.WriteHeader(statusCode)
}

func (w *responseWriter) Flush() {
	if VERBOSE_LOGGING && w.statusCode >= 400 {
		correlationId := w.ctx.Value(CORRELATION_ID).(string)
		fmt.Fprintf(os.Stderr, "---begin: %s---\n", correlationId)
		fmt.Fprintln(os.Stderr, w.buf.String())
		fmt.Fprintf(os.Stderr, "---end: %s---\n", correlationId)
		w.buf.Reset()
	}

	if f, ok := w.parent.(http.Flusher); ok {
		f.Flush()
	}
}
