package logging

import (
	"context"

	"github.com/google/uuid"
)

type key string

const CORRELATION_ID = key("correlationId")
const SENSITIVE_MARKER = key("sensitiveMarker")

type sensitiveMarker struct {
	marked bool
}

func WithCorrelationId(
	ctx context.Context,
	correlationId string,
) context.Context {
	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	return context.WithValue(ctx, CORRELATION_ID, correlationId)
}

func WithSensitiveMarker(ctx context.Context) context.Context {
	return context.WithValue(ctx, SENSITIVE_MARKER, &sensitiveMarker{marked: false})
}

func MarkSensitive(ctx context.Context) {
	m := ctx.Value(SENSITIVE_MARKER)
	if m != nil {
		m.(*sensitiveMarker).marked = true
	}
}

func IsMarkedSensitive(ctx context.Context) bool {
	m := ctx.Value(SENSITIVE_MARKER)
	if m == nil {
		return false
	}
	return m.(*sensitiveMarker).marked
}
