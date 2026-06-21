package cluster

import (
	"context"
	"errors"

	"encoding/json"
	"fmt"

	"link-society.com/flowg/internal/storage/changefeed"
)

var (
	errInvalidNotification = errors.New("invalid notification")
)

const (
	writeNotificationTag byte = iota + 1
)

type notification interface {
	Marshal() []byte
	Handle(ctx context.Context, delegate *delegate) error
}

type writeNotification struct {
	Namespace string              `json:"namespace"`
	Records   []changefeed.Record `json:"records"`
}

var _ notification = (*writeNotification)(nil)

func (n *writeNotification) Marshal() []byte {
	buf, err := json.Marshal(n)
	if err != nil {
		return nil
	}
	return append([]byte{writeNotificationTag}, buf...)
}

func (n *writeNotification) Handle(ctx context.Context, d *delegate) error {
	s, ok := d.storages[n.Namespace]
	if !ok {
		return fmt.Errorf("unknown namespace %q in write notification", n.Namespace)
	}
	return s.ApplyReplicated(ctx, n.Records)
}

func parseNotification(data []byte) (notification, error) {
	if len(data) == 0 {
		return nil, errInvalidNotification
	}

	switch data[0] {
	case writeNotificationTag:
		var n writeNotification
		if err := json.Unmarshal(data[1:], &n); err != nil {
			return nil, err
		}
		return &n, nil

	default:
		return nil, errInvalidNotification
	}
}
