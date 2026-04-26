package cluster

import (
	"context"
	"log/slog"
	"time"

	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"
	"go.uber.org/fx"

	"link-society.com/flowg/internal/utils/fxproviders"

	"link-society.com/flowg/internal/storage"
	"link-society.com/flowg/internal/storage/auth"
	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"github.com/hashicorp/memberlist"
)

type Manager interface {
	actor.Actor

	HttpHandler() http.Handler
}

type ManagerOptions struct {
	NodeID string
	Cookie string

	ClusterFormationStrategy ClusterFormationStrategy
	ClusterStateDir          string
}

type managerImpl struct {
	actor.Actor

	handler http.Handler
	notifyM actor.Mailbox[notification]
}

var _ Manager = (*managerImpl)(nil)

func NewManager(opts ManagerOptions) fx.Option {
	return fx.Module(
		"cluster.manager",
		clusterstate.NewStorage(func() clusterstate.Options {
			clusterStateOpts := clusterstate.DefaultOptions()
			clusterStateOpts.Directory = opts.ClusterStateDir
			return clusterStateOpts
		}()),
		fxproviders.ProvideMailbox[notification](),
		fxproviders.ProvideMailbox[net.Conn](),
		fxproviders.ProvideMailbox[*memberlist.Packet](),
		fxproviders.ProvideMailbox[*ClusterJoinNode](),
		fxproviders.ProvideMailbox[*syncRequest](),
		fxproviders.ProvideActor[*syncActor](
			func(d struct {
				fx.In

				SyncRequestM actor.Mailbox[*syncRequest]

				AuthStorage         auth.Storage
				ConfigStorage       config.Storage
				LogStorage          log.Storage
				ClusterStateStorage clusterstate.Storage
			}) *syncActor {
				return &syncActor{
					Actor: actor.New(&syncWorker{
						logger: slog.Default().With(slog.String("channel", "cluster.sync")),

						requestM: d.SyncRequestM,
						cookie:   opts.Cookie,

						storages: map[string]storage.Streamable{
							"auth":   d.AuthStorage,
							"config": d.ConfigStorage,
							"log":    d.LogStorage,
						},
						clusterStateStorage: d.ClusterStateStorage,
					}),
				}
			},
		),
		fxproviders.ProvideActor[*clusterFormationController](
			func(
				joinM actor.Mailbox[*ClusterJoinNode],
				listener *Listener,
			) *clusterFormationController {
				return &clusterFormationController{
					Actor: actor.New(&clusterFormationControllerWorker{
						logger:   slog.Default().With(slog.String("channel", "cluster.formation")),
						joinM:    joinM,
						resolver: func() (*url.URL, error) { return listener.ResolveLocalEndpoint() },
						strategy: opts.ClusterFormationStrategy,
					}),
				}
			},
		),
		fx.Provide(func(d struct {
			fx.In

			Listener            *Listener
			SyncRequestM        actor.Mailbox[*syncRequest]
			ClusterStateStorage clusterstate.Storage
		}) (*delegate, error) {
			var err error

			localEndpoint, err := d.Listener.ResolveLocalEndpoint()
			if err != nil {
				return nil, err
			}

			logger := slog.Default().With(
				slog.String("channel", "cluster.gossip"),
				slog.String("cluster.local.node", opts.NodeID),
				slog.String("cluster.local.endpoint", localEndpoint.String()),
			)

			return &delegate{
				logger: logger,

				localNodeID:   opts.NodeID,
				localEndpoint: localEndpoint,
				endpoints:     newEndpointCache(),

				notifyC: make(chan notification, 1000),

				clusterStateStorage: d.ClusterStateStorage,
				syncRequestM:        d.SyncRequestM,
			}, nil
		}),
		fx.Provide(func(
			d struct {
				fx.In

				Delegate *delegate
				ConnM    actor.Mailbox[net.Conn]
				PacketM  actor.Mailbox[*memberlist.Packet]

				AuthStorage         auth.Storage
				ConfigStorage       config.Storage
				LogStorage          log.Storage
				ClusterStateStorage clusterstate.Storage
			},
		) *httpTransport {
			return &httpTransport{
				delegate: d.Delegate,
				cookie:   opts.Cookie,

				connM:   d.ConnM,
				packetM: d.PacketM,

				storages: map[string]storage.Streamable{
					"auth":   d.AuthStorage,
					"config": d.ConfigStorage,
					"log":    d.LogStorage,
				},
				clusterStateStorage: d.ClusterStateStorage,
			}
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			delegate *delegate,
			transport *httpTransport,
		) (*memberlist.Memberlist, error) {
			mlistConfig := memberlist.DefaultLocalConfig()
			mlistConfig.Name = opts.NodeID
			mlistConfig.RequireNodeNames = true
			mlistConfig.Transport = transport
			mlistConfig.Delegate = delegate
			mlistConfig.Events = delegate
			mlistConfig.PushPullInterval = time.Second
			mlistConfig.Logger = newMemberlistLogger(delegate.logger)

			mlist, err := memberlist.Create(mlistConfig)
			if err != nil {
				return nil, err
			}

			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					if err := mlist.Leave(5 * time.Second); err != nil {
						return err
					}

					return mlist.Shutdown()
				},
			})

			return mlist, nil
		}),
		fxproviders.ProvideActor[Manager](
			func(d struct {
				fx.In

				JoinM   actor.Mailbox[*ClusterJoinNode]
				NotifyM actor.Mailbox[notification]

				Delegate   *delegate
				Memberlist *memberlist.Memberlist
				Transport  *httpTransport
			}) Manager {
				worker := actor.NewWorker(func(ctx actor.Context) actor.WorkerStatus {
					select {
					case <-ctx.Done():
						return actor.WorkerEnd

					case msg, ok := <-d.NotifyM.ReceiveC():
						if !ok {
							return actor.WorkerEnd
						}

						payload := msg.Marshal()
						for _, node := range d.Memberlist.Members() {
							if node.Name != d.Delegate.localNodeID {
								go d.Memberlist.SendReliable(node, payload)
							}
						}

						return actor.WorkerContinue

					case msg, ok := <-d.Delegate.notifyC:
						if !ok {
							return actor.WorkerEnd
						}

						msg.Handle(ctx, d.Delegate)

						return actor.WorkerContinue

					case joinNode, ok := <-d.JoinM.ReceiveC():
						if !ok {
							return actor.WorkerEnd
						}

						d.Delegate.endpoints.Set(joinNode.JoinNodeID, joinNode.JoinNodeEndpoint)
						_, err := d.Memberlist.Join([]string{joinNode.Address()})
						if err != nil {
							d.Delegate.logger.Error(
								"failed to join cluster node",
								slog.String("cluster.join.node", joinNode.JoinNodeID),
								slog.String("cluster.join.address", joinNode.Address()),
								slog.String("error", err.Error()),
							)
							return actor.WorkerEnd
						}

						return actor.WorkerContinue
					}
				})

				return &managerImpl{
					Actor:   actor.New(worker),
					handler: d.Transport,
					notifyM: d.NotifyM,
				}
			},
		),
		fx.Invoke(func(_ struct {
			S *syncActor
			C *clusterFormationController
		}) {
			// No-op, just to force the creation of all components
		}),
	)
}

func (m *managerImpl) HttpHandler() http.Handler {
	return m.handler
}
