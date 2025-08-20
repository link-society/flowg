package cluster

import (
	"context"
	"log/slog"
	"time"

	"encoding/json"

	"net/url"

	"github.com/hashicorp/memberlist"

	"link-society.com/flowg/internal/utils/kvstore"
)

type delegate struct {
	logger *slog.Logger

	localNodeID   string
	localEndpoint *url.URL
	endpoints     map[string]*url.URL

	clusterStateStorage kvstore.Storage
	syncPool            *syncPool
}

var _ memberlist.Delegate = (*delegate)(nil)

func (d *delegate) NodeMeta(int) []byte {
	return []byte(d.localEndpoint.String())
}

func (d *delegate) NotifyMsg([]byte) {
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{}
}

func (d *delegate) LocalState(join bool) []byte {
	endpointNames := make([]string, 0, len(d.endpoints))
	for endpointName := range d.endpoints {
		endpointNames = append(endpointNames, endpointName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	localState, err := fetchLocalState(ctx, d.clusterStateStorage, d.localNodeID, endpointNames)
	if err != nil {
		d.logger.Error(
			"failed to get local state",
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
	var remoteState nodeState
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
				"failed to get last sync",
				slog.String("cluster.remote.node", remoteState.NodeID),
			)

			return
		}

		worker, ok := d.syncPool.workers[remoteState.NodeID]
		if !ok {
			d.logger.Error(
				"no sync worker for node",
				slog.String("cluster.remote.node", remoteState.NodeID),
			)

			return
		}

		err := worker.mbox.Send(context.Background(), lastSync)
		if err != nil {
			d.logger.Error(
				"failed to notify sync worker",
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

	d.endpoints[node.Name] = endpointUrl
	d.syncPool.AddWorker(node.Name, endpointUrl)
}

func (d *delegate) NotifyLeave(node *memberlist.Node) {
	d.logger.Info(
		"remote node left",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", string(node.Meta)),
	)

	d.syncPool.RemoveWorker(node.Name)
	delete(d.endpoints, node.Name)
}

func (d *delegate) NotifyUpdate(node *memberlist.Node) {
	d.logger.Info(
		"remote node updated",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", string(node.Meta)),
	)

	d.NotifyJoin(node)
}
