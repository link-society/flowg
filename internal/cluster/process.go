package cluster

import (
	"errors"
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
	joinM   actor.MailboxReceiver[*ClusterJoinNode]

	delegate    *delegate
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

	p.delegate = &delegate{
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

	transport := &httpTransport{
		delegate: p.delegate,
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
	p.mlistConfig.Delegate = p.delegate
	p.mlistConfig.Events = p.delegate
	p.mlistConfig.PushPullInterval = time.Second
	p.mlistConfig.Logger = newMemberlistLogger(logger)

	p.mlist, err = memberlist.Create(p.mlistConfig)
	if err != nil {
		return proctree.Terminate(err)
	}

	p.httpHandler = transport

	return proctree.Continue()
}

func (p *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case joinNode, ok := <-p.joinM.ReceiveC():
		if !ok {
			return proctree.Terminate(nil)
		}

		p.delegate.endpoints[joinNode.JoinNodeID] = joinNode.JoinNodeEndpoint
		_, err := p.mlist.Join([]string{joinNode.Address()})
		if err != nil {
			return proctree.Terminate(err)
		}

		return proctree.Continue()
	}
}

func (p *procHandler) Terminate(ctx actor.Context, parentErr error) error {
	if p.mlistConfig != nil && p.mlistConfig.Delegate != nil {
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

	return parentErr
}
