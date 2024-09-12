package logging

import (
	"context"

	"github.com/google/uuid"
)

type key string

const CORRELATION_ID = key("correlationId")

func WithCorrelationId(
	ctx context.Context,
	correlationId string,
) context.Context {
	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	return context.WithValue(ctx, CORRELATION_ID, correlationId)
}
