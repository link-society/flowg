package cluster

import (
	"context"
	"errors"
)

var (
	errInvalidNotification = errors.New("invalid notification")
)

const (
	_ = iota
)

type notification interface {
	Marshal() []byte
	Handle(ctx context.Context, delegate *delegate) error
}

func parseNotification(data []byte) (notification, error) {
	if len(data) == 0 {
		return nil, errInvalidNotification
	}

	switch data[0] {
	default:
		return nil, errInvalidNotification
	}
}
