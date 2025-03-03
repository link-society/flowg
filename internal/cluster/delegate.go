package cluster

import (
	"log/slog"

	"net/url"

	"github.com/hashicorp/memberlist"
)

type delegate struct {
	logger *slog.Logger

	localEndpoint *url.URL
	endpoints     map[string]*url.URL
}

func (d *delegate) NodeMeta(int) []byte {
	return []byte(d.localEndpoint.String())
}

func (d *delegate) NotifyMsg([]byte) {
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return [][]byte{}
}

func (d *delegate) LocalState(join bool) []byte {
	return []byte{}
}

func (d *delegate) MergeRemoteState([]byte, bool) {
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
