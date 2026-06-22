package cluster

import (
	"context"
	"time"

	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

func namespaceStaleness(
	ctx context.Context,
	store clusterstate.Storage,
	namespace string,
	threshold time.Duration,
) (stale bool, err error) {
	liveness, err := store.GetLiveness(ctx, namespace)
	if err != nil {
		return false, err
	}
	if liveness == 0 {
		return true, nil
	}
	return time.Since(time.Unix(0, liveness)) >= threshold, nil
}
