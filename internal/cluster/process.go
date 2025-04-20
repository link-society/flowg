package cluster

import (
	"errors"
	"fmt"
	"log/slog"

	"time"

	"net"
	"net/http"
	"net/url"

	"github.com/hashicorp/memberlist"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	opts *ManagerOptions

	connM   actor.Mailbox[net.Conn]
	packetM actor.Mailbox[*memberlist.Packet]

	mlistConfig *memberlist.Config
	mlist       *memberlist.Memberlist

	httpHandler http.Handler
}

var _ proctree.ProcessHandler = (*procHandler)(nil)

func (p *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	var err error

	localEndpoint, err := p.opts.LocalEndpointResolver()
	if err != nil {
		return proctree.Terminate(err)
	}

	logger := slog.Default().With(
		slog.String("channel", "cluster.gossip"),
		slog.String("cluster.local.node", p.opts.NodeID),
		slog.String("cluster.local.endpoint", localEndpoint.String()),
	)

	d := &delegate{
		logger: logger,

		localEndpoint: localEndpoint,
		endpoints:     make(map[string]*url.URL),
	}

	if p.opts.ClusterJoinNode.JoinNodeID != "" && p.opts.ClusterJoinNode.JoinNodeEndpoint != nil {
		d.endpoints[p.opts.ClusterJoinNode.JoinNodeID] = p.opts.ClusterJoinNode.JoinNodeEndpoint
	}

	transport := &httpTransport{
		delegate: d,
		cookie:   p.opts.Cookie,

		connM:   p.connM,
		packetM: p.packetM,
	}

	p.mlistConfig = memberlist.DefaultLocalConfig()
	p.mlistConfig.Name = p.opts.NodeID
	p.mlistConfig.RequireNodeNames = true
	p.mlistConfig.Transport = transport
	p.mlistConfig.Delegate = d
	p.mlistConfig.Events = d
	p.mlistConfig.Logger = newMemberlistLogger(logger)

	p.mlist, err = memberlist.Create(p.mlistConfig)
	if err != nil {
		return proctree.Terminate(err)
	}

	/* If automatic cluster formation is enabled,
	then ConsulService already discovered nodes and set them in ClusterNodeNode.
	ConsulService starts before MagementServer therefore ManagementServer need not
	wait here for ConsulService to finish discovering nodes */
	if p.opts.ClusterJoinNode.JoinNodeID != "" && p.opts.ClusterJoinNode.JoinNodeEndpoint != nil {
		joinAddr := fmt.Sprintf("%s/%s", p.opts.ClusterJoinNode.JoinNodeID, p.opts.ClusterJoinNode.JoinNodeEndpoint.Host)
		_, err = p.mlist.Join([]string{joinAddr})
		if err != nil {
			if !p.opts.AutomaticClusterFormation {
				/* Terminate the process only when AutomaticClusterFormation is disabled
				because otherwise the first node in the cluster will never know other nodes to connect to
				therefore we cannot terminate the process for the very first node in the cluster */
				logger.ErrorContext(ctx,
					"memberlist join failed",
					slog.Any("error", err),
				)
				return proctree.Terminate(err)
			} else {
				logger.WarnContext(ctx,
					"memberlist join failed",
					slog.Any("error", err),
				)
			}
		}
	}

	p.httpHandler = transport

	return proctree.Continue()
}

func (p *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (p *procHandler) Terminate(ctx actor.Context, parentErr error) error {
	if p.mlist != nil {
		if err := p.mlist.Leave(5 * time.Second); err != nil {
			return errors.Join(parentErr, err)
		}

		if err := p.mlist.Shutdown(); err != nil {
			return errors.Join(parentErr, err)
		}
	}

	return parentErr
}
