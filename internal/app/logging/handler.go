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

func newHandler(w io.Writer, opts *slog.HandlerOptions) *handler {
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
