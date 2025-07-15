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
	joinNode    *ClusterJoinNode

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

		localNodeID:   p.opts.NodeID,
		localEndpoint: localEndpoint,
		endpoints:     make(map[string]*url.URL),

		clusterStateStorage: p.opts.ClusterStateStorage,

		syncPool: &syncPool{
			logger: slog.Default().With(
				slog.String("channel", "cluster.replication"),
				slog.String("cluster.local.node", p.opts.NodeID),
				slog.String("cluster.local.endpoint", localEndpoint.String()),
			),

			nodeID: p.opts.NodeID,
			cookie: p.opts.Cookie,

			authStorage:   p.opts.AuthStorage,
			configStorage: p.opts.ConfigStorage,
			logStorage:    p.opts.LogStorage,

			workers: make(map[string]*syncActor),
		},
	}

	p.joinNode, err = p.opts.ClusterFormationStrategy.Join(ctx, p.opts.LocalEndpointResolver)
	if err != nil {
		return proctree.Terminate(fmt.Errorf("failed to join cluster: %w", err))
	}

	if !p.joinNode.IsEmpty() {
		d.endpoints[p.joinNode.JoinNodeID] = p.joinNode.JoinNodeEndpoint
	}

	transport := &httpTransport{
		delegate: d,
		cookie:   p.opts.Cookie,

		connM:   p.connM,
		packetM: p.packetM,

		authStorage:   p.opts.AuthStorage,
		configStorage: p.opts.ConfigStorage,
		logStorage:    p.opts.LogStorage,
	}

	p.mlistConfig = memberlist.DefaultLocalConfig()
	p.mlistConfig.Name = p.opts.NodeID
	p.mlistConfig.RequireNodeNames = true
	p.mlistConfig.Transport = transport
	p.mlistConfig.Delegate = d
	p.mlistConfig.Events = d
	p.mlistConfig.PushPullInterval = time.Second
	p.mlistConfig.Logger = newMemberlistLogger(logger)

	p.mlist, err = memberlist.Create(p.mlistConfig)
	if err != nil {
		return proctree.Terminate(err)
	}

	if !p.joinNode.IsEmpty() {
		_, err = p.mlist.Join([]string{p.joinNode.Address()})
		if err != nil {
			logger.ErrorContext(
				ctx,
				"memberlist join failed",
				slog.Any("error", err),
			)
			return proctree.Terminate(err)
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
	if p.mlistConfig.Delegate != nil {
		d := p.mlistConfig.Delegate.(*delegate)
		d.syncPool.RemoveAll()
	}

	if p.mlist != nil {
		if err := p.mlist.Leave(5 * time.Second); err != nil {
			return errors.Join(parentErr, err)
		}

		if err := p.mlist.Shutdown(); err != nil {
			return errors.Join(parentErr, err)
		}
	}

	if p.joinNode != nil {
		if err := p.opts.ClusterFormationStrategy.Leave(ctx, p.joinNode); err != nil {
			return errors.Join(parentErr, fmt.Errorf("failed to leave cluster: %w", err))
		}
	}

	return parentErr
}
