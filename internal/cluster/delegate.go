package cluster

import (
	"log/slog"

	"net/url"

	"github.com/hashicorp/memberlist"
)

type delegate struct {
	logger *slog.Logger

	localNodeID   string
	localEndpoint *url.URL
	endpoints     map[string]*url.URL

	notifyC chan notification
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
	return []byte{}
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
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
}

func (d *delegate) NotifyLeave(node *memberlist.Node) {
	d.logger.Info(
		"remote node left",
		slog.String("cluster.remote.node", node.Name),
		slog.String("cluster.remote.endpoint", string(node.Meta)),
	)

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
