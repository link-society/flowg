package cluster

import (
	"context"
	"errors"

	"github.com/hashicorp/memberlist"
)

var (
	errInvalidBroadcastMessage = errors.New("invalid broadcast message")
)

const (
	_                                                      = iota
	BROADCAST_MESSAGE_TYPE_INVALIDATE_PIPELINE_BUILD_CACHE = iota
)

type broadcastMessage interface {
	memberlist.Broadcast

	Handle(ctx context.Context, delegate *delegate) error
}

type invalidatePipelineBuildCache struct {
	pipelineName string // if empty, then all
}

var _ broadcastMessage = (*invalidatePipelineBuildCache)(nil)

func (msg *invalidatePipelineBuildCache) Handle(ctx context.Context, delegate *delegate) error {
	if msg.pipelineName == "" {
		return delegate.pipelineRunner.InvalidateAllCachedBuilds(ctx)
	} else {
		return delegate.pipelineRunner.InvalidateCachedBuild(ctx, msg.pipelineName)
	}
}

func (*invalidatePipelineBuildCache) Invalidates(memberlist.Broadcast) bool {
	return false
}

func (msg *invalidatePipelineBuildCache) Message() []byte {
	var payload []byte
	payload = append(payload, BROADCAST_MESSAGE_TYPE_INVALIDATE_PIPELINE_BUILD_CACHE)
	payload = append(payload, msg.pipelineName...)
	return payload
}

func (*invalidatePipelineBuildCache) Finished() {
	// No-op
}

func parseBroadcastMessage(data []byte) (broadcastMessage, error) {
	if len(data) == 0 {
		return nil, errInvalidBroadcastMessage
	}

	switch data[0] {
	case BROADCAST_MESSAGE_TYPE_INVALIDATE_PIPELINE_BUILD_CACHE:
		return &invalidatePipelineBuildCache{
			pipelineName: string(data[1:]),
		}, nil
	default:
		return nil, errInvalidBroadcastMessage
	}
}
