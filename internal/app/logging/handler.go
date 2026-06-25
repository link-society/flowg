package logging

import (
	"context"
	"io"
	"log/slog"
)

type handler struct {
	parent slog.Handler
}

var _ slog.Handler = (*handler)(nil)

// VERBOSE_LOGGING controls whether the access-log middleware buffers response
// bodies and dumps them to stderr for failed requests. It is toggled by the
// server's logging setup and read by the API access-log middleware.
var VERBOSE_LOGGING = false

// NewHandler builds a slog.Handler that writes text-formatted records to w and
// enriches each record with the correlation id carried by the context (see
// CORRELATION_ID).
func NewHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return &handler{
		parent: slog.NewTextHandler(w, opts),
	}
}

func (h *handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.parent.Enabled(ctx, lvl)
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	if v, ok := ctx.Value(CORRELATION_ID).(string); ok {
		record.AddAttrs(slog.String("correlation_id", v))
	}

	return h.parent.Handle(ctx, record)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{
		parent: h.parent.WithAttrs(attrs),
	}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{
		parent: h.parent.WithGroup(name),
	}
}
