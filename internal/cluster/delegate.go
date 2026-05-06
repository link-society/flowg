package cluster

import (
	"context"
	"log/slog"

	"encoding/json"
	"time"

	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"github.com/hashicorp/memberlist"

	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type delegate struct {
	logger *slog.Logger

	localNodeID   string
	localEndpoint *url.URL
	endpoints     *endpointCache

	notifyC chan notification

	clusterStateStorage clusterstate.Storage
	syncRequestM        actor.MailboxSender[*syncRequest]
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
	} else {
		d.notifyC <- msg
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

	if remoteState.NodeID != d.localNodeID {
		lastSync, ok := remoteState.LastSync[d.localNodeID]
		if !ok {
			d.logger.Error(
				"remote state does not contain sync information for local node",
				slog.String("cluster.remote.node", remoteState.NodeID),
			)
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

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		req := &syncRequest{
			remoteNodeID:   remoteState.NodeID,
			remoteEndpoint: remoteEndpoint,
			lastSync:       lastSync,
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
