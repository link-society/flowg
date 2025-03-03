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
	nodeID        string
	localEndpoint *url.URL

	joinNodeID       string
	joinNodeEndpoint *url.URL

	connM   actor.Mailbox[net.Conn]
	packetM actor.Mailbox[*memberlist.Packet]

	mlistConfig *memberlist.Config
	mlist       *memberlist.Memberlist

	httpHandler http.Handler
}

func (p *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	var err error

	logger := slog.Default().With(
		slog.String("channel", "cluster.gossip"),
		slog.String("cluster.local.node", p.nodeID),
		slog.String("cluster.local.endpoint", p.localEndpoint.String()),
	)

	d := &delegate{
		logger: logger,

		localEndpoint: p.localEndpoint,
		endpoints:     make(map[string]*url.URL),
	}

	if p.joinNodeID != "" && p.joinNodeEndpoint != nil {
		d.endpoints[p.joinNodeID] = p.joinNodeEndpoint
	}

	transport := &httpTransport{
		delegate: d,
		connM:    p.connM,
		packetM:  p.packetM,
	}

	p.mlistConfig = memberlist.DefaultLocalConfig()
	p.mlistConfig.Name = p.nodeID
	p.mlistConfig.RequireNodeNames = true
	p.mlistConfig.Transport = transport
	p.mlistConfig.Delegate = d
	p.mlistConfig.Events = d
	p.mlistConfig.Logger = newMemberlistLogger(logger)

	p.mlist, err = memberlist.Create(p.mlistConfig)
	if err != nil {
		return proctree.Terminate(err)
	}

	if p.joinNodeID != "" && p.joinNodeEndpoint != nil {
		joinAddr := fmt.Sprintf("%s/%s", p.joinNodeID, p.joinNodeEndpoint.Host)
		_, err = p.mlist.Join([]string{joinAddr})
		if err != nil {
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
