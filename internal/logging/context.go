package logging

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type key string

const CORRELATION_ID = key("correlationId")

type correlatedContext struct {
	parent        context.Context
	correlationId string
}

func (c *correlatedContext) Deadline() (deadline time.Time, ok bool) {
	return c.parent.Deadline()
}

func (c *correlatedContext) Done() <-chan struct{} {
	return c.parent.Done()
}

func (c *correlatedContext) Err() error {
	return c.parent.Err()
}

func (c *correlatedContext) Value(key any) any {
	if key == CORRELATION_ID {
		return c.correlationId
	}
	return c.parent.Value(key)
}

func WithCorrelationId(
	ctx context.Context,
	correlationId string,
) context.Context {
	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	return &correlatedContext{
		parent:        ctx,
		correlationId: correlationId,
	}
}
