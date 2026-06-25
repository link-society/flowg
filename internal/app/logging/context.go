package logging

import (
	"context"

	"github.com/google/uuid"
)

type key string

// CORRELATION_ID is the context key under which the request correlation id is
// stored. The slog handler reads it to attach a `correlation_id` attribute to
// every log record emitted within the request scope.
const CORRELATION_ID = key("correlationId")

// SENSITIVE_MARKER is the context key under which the sensitive-request marker
// is stored.
const SENSITIVE_MARKER = key("sensitiveMarker")

type sensitiveMarker struct {
	marked bool
}

// WithCorrelationId returns a copy of ctx carrying the given correlation id. If
// correlationId is empty, a new random UUID is generated.
func WithCorrelationId(
	ctx context.Context,
	correlationId string,
) context.Context {
	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	return context.WithValue(ctx, CORRELATION_ID, correlationId)
}

// WithSensitiveMarker returns a copy of ctx carrying a sensitive marker that is
// initially unset. Handlers may call MarkSensitive to flag the request as
// carrying sensitive data.
func WithSensitiveMarker(ctx context.Context) context.Context {
	return context.WithValue(ctx, SENSITIVE_MARKER, &sensitiveMarker{marked: false})
}

// MarkSensitive flags the request associated with ctx as carrying sensitive
// data. It is a no-op if ctx does not carry a sensitive marker.
func MarkSensitive(ctx context.Context) {
	m := ctx.Value(SENSITIVE_MARKER)
	if m != nil {
		m.(*sensitiveMarker).marked = true
	}
}

// IsMarkedSensitive reports whether the request associated with ctx has been
// flagged as carrying sensitive data.
func IsMarkedSensitive(ctx context.Context) bool {
	m := ctx.Value(SENSITIVE_MARKER)
	if m == nil {
		return false
	}
	return m.(*sensitiveMarker).marked
}
