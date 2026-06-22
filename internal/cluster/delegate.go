package cluster

import (
	"context"
	"log/slog"

	"encoding/json"
	"time"

	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"github.com/hashicorp/memberlist"

	"link-society.com/flowg/internal/storage"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type delegate struct {
	logger *slog.Logger

	localNodeID   string
	localEndpoint *url.URL
	endpoints     *endpointCache

	notificationConsumerM actor.MailboxSender[notification]

	bootstrapThreshold time.Duration

	clusterStateStorage clusterstate.Storage
	syncRequestM        actor.MailboxSender[*syncRequest]
	watermarks          *watermarkCache
	storages            map[string]storage.Streamable
}

var _ memberlist.Delegate = (*delegate)(nil)

func (d *delegate) NodeMeta(int) []byte {
	return []byte(d.localEndpoint.String())
}

func (d *delegate) NotifyMsg(msg []byte) {
	if msg, err := parseNotification(msg); err != nil {
		d.logger.Error(
			"failed to parse notification message",
			slog.String("error", err.Error()),
		)
	} else if err := d.notificationConsumerM.Send(context.Background(), msg); err != nil {
		d.logger.Error(
			"failed to enqueue notification message",
			slog.String("error", err.Error()),
		)
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{}
}

func (d *delegate) LocalState(join bool) []byte {
	endpointNames := []string{}

	for endpointName := range d.endpoints.All() {
		endpointNames = append(endpointNames, endpointName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	localState, err := d.clusterStateStorage.FetchLocalState(
		ctx,
		d.localNodeID,
		endpointNames,
	)
	if err != nil {
		d.logger.Error(
			"failed to fetch local state",
			slog.String("error", err.Error()),
		)
		return []byte{}
	}

	localState.Healthy = make(map[string]bool, len(d.storages))
	for namespace := range d.storages {
		stale, err := namespaceStaleness(ctx, d.clusterStateStorage, namespace, d.bootstrapThreshold)
		if err != nil {
			d.logger.Error(
				"failed to evaluate namespace staleness",
				slog.String("cluster.replication.namespace", namespace),
				slog.String("error", err.Error()),
			)
		}
		localState.Healthy[namespace] = !stale
	}

	buf, err := json.Marshal(localState)
	if err != nil {
		d.logger.Error(
			"failed to encode local state",
			slog.String("error", err.Error()),
		)
		return []byte{}
	}

	return buf
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	var remoteState clusterstate.NodeState
	if err := json.Unmarshal(buf, &remoteState); err != nil {
		d.logger.Error(
			"failed to decode remote state",
			slog.String("error", err.Error()),
		)
		return
	}

	if remoteState.NodeID == d.localNodeID {
		return
	}

	remoteEndpoint, ok := d.endpoints.Get(remoteState.NodeID)
	if !ok {
		d.logger.Error(
			"remote endpoint not found",
			slog.String("cluster.remote.node", remoteState.NodeID),
		)
		return
	}

	knownSince := make(map[string]uint64)
	for _, syncState := range remoteState.LastSync[d.localNodeID] {
		knownSince[syncState.Namespace] = syncState.Since
		if d.watermarks != nil {
			d.watermarks.observe(remoteState.NodeID, syncState.Namespace, syncState.Since)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	now := time.Now().UnixNano()

	lastSync := make([]clusterstate.NamespaceSyncState, 0, len(d.storages))
	bootstrap := make([]string, 0)

	for namespace, store := range d.storages {
		stale, err := namespaceStaleness(ctx, d.clusterStateStorage, namespace, d.bootstrapThreshold)
		if err != nil {
			d.logger.Error(
				"failed to evaluate namespace staleness",
				slog.String("cluster.replication.namespace", namespace),
				slog.String("error", err.Error()),
			)
		}

		if stale {
			switch {
			case remoteState.Healthy[namespace]:
				bootstrap = append(bootstrap, namespace)

			case d.shouldSelfHeal(ctx, namespace):
				if err := d.clusterStateStorage.SetLiveness(ctx, namespace, now); err != nil {
					d.logger.Error(
						"failed to self-heal liveness",
						slog.String("cluster.replication.namespace", namespace),
						slog.String("error", err.Error()),
					)
				} else {
					d.logger.Warn(
						"no healthy bootstrap source found; resuming replication",
						slog.String("cluster.replication.namespace", namespace),
					)
				}
			}
			continue
		}

		if err := d.clusterStateStorage.SetLiveness(ctx, namespace, now); err != nil {
			d.logger.Error(
				"failed to update liveness",
				slog.String("cluster.replication.namespace", namespace),
				slog.String("error", err.Error()),
			)
		}

		latest, err := store.LatestVersion(ctx)
		if err != nil {
			d.logger.Error(
				"failed to fetch latest version",
				slog.String("cluster.replication.namespace", namespace),
				slog.String("error", err.Error()),
			)
			continue
		}

		if knownSince[namespace] >= latest {
			continue
		}

		lastSync = append(lastSync, clusterstate.NamespaceSyncState{
			Namespace: namespace,
			Since:     knownSince[namespace],
		})
	}

	if len(lastSync) == 0 && len(bootstrap) == 0 {
		return
	}

	req := &syncRequest{
		remoteNodeID:   remoteState.NodeID,
		remoteEndpoint: remoteEndpoint,
		lastSync:       lastSync,
		bootstrap:      bootstrap,
	}
	if err := d.syncRequestM.Send(ctx, req); err != nil {
		d.logger.Error(
			"failed to send sync request",
			slog.String("cluster.remote.node", remoteState.NodeID),
			slog.String("error", err.Error()),
		)
		return
	}
}

func (d *delegate) shouldSelfHeal(ctx context.Context, namespace string) bool {
	liveness, err := d.clusterStateStorage.GetLiveness(ctx, namespace)
	if err != nil {
		d.logger.Error(
			"failed to read liveness",
			slog.String("cluster.replication.namespace", namespace),
			slog.String("error", err.Error()),
		)
		return false
	}
	if liveness == 0 {
		return false
	}
	return time.Since(time.Unix(0, liveness)) >= 2*d.bootstrapThreshold
}

func (d *delegate) NotifyJoin(node *memberlist.Node) {
	endpoint := string(node.Meta)

	d.logger.Info(
		"remote node joined",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", endpoint),
	)

	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		d.logger.Error(
			"failed to parse endpoint",
			slog.String("cluster.remote.node", node.Name),
			slog.String("cluster.remote.endpoint", endpoint),
			slog.String("error", err.Error()),
		)
		return
	}

	d.endpoints.Set(node.Name, endpointUrl)
}

func (d *delegate) NotifyLeave(node *memberlist.Node) {
	d.logger.Info(
		"remote node left",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", string(node.Meta)),
	)

	d.endpoints.Delete(node.Name)
}

func (d *delegate) NotifyUpdate(node *memberlist.Node) {
	d.logger.Info(
		"remote node updated",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", string(node.Meta)),
	)

	d.NotifyJoin(node)
}
