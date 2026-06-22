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
	"link-society.com/flowg/internal/storage/changefeed"
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

	TombstoneGracePeriod time.Duration
}

type managerImpl struct {
	actor.Actor

	handler                http.Handler
	notificationPublisherM actor.Mailbox[notification]
}

var _ Manager = (*managerImpl)(nil)

type notificationMailboxes struct {
	actor.Actor

	publisher actor.Mailbox[notification]
	consumer  actor.Mailbox[notification]
}

func newNotificationMailboxes() *notificationMailboxes {
	publisher := actor.NewMailbox[notification]()
	consumer := actor.NewMailbox[notification]()

	return &notificationMailboxes{
		Actor:     actor.Combine(publisher, consumer).Build(),
		publisher: publisher,
		consumer:  consumer,
	}
}

func NewManager(opts ManagerOptions) fx.Option {
	bootstrapThreshold := opts.TombstoneGracePeriod / 2

	return fx.Module(
		"cluster.manager",
		clusterstate.NewStorage(func() clusterstate.Options {
			clusterStateOpts := clusterstate.DefaultOptions()
			clusterStateOpts.Directory = opts.ClusterStateDir
			return clusterStateOpts
		}()),
		fxproviders.ProvideMailbox[net.Conn](),
		fxproviders.ProvideMailbox[*memberlist.Packet](),
		fxproviders.ProvideMailbox[*ClusterJoinNode](),
		fxproviders.ProvideMailbox[*syncRequest](),
		fxproviders.ProvideActor[*notificationMailboxes](newNotificationMailboxes),
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

						localNodeID:        opts.NodeID,
						requestM:           d.SyncRequestM,
						cookie:             opts.Cookie,
						bootstrapThreshold: bootstrapThreshold,

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

			Listener              *Listener
			SyncRequestM          actor.Mailbox[*syncRequest]
			ClusterStateStorage   clusterstate.Storage
			NotificationMailboxes *notificationMailboxes

			AuthStorage   auth.Storage
			ConfigStorage config.Storage
			LogStorage    log.Storage
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

				notificationConsumerM: d.NotificationMailboxes.consumer,

				bootstrapThreshold: bootstrapThreshold,

				clusterStateStorage: d.ClusterStateStorage,
				syncRequestM:        d.SyncRequestM,
				storages: map[string]storage.Streamable{
					"auth":   d.AuthStorage,
					"config": d.ConfigStorage,
					"log":    d.LogStorage,
				},
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

				bootstrapThreshold: bootstrapThreshold,

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
			mlistConfig.PushPullInterval = 10 * time.Second
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

				JoinM                 actor.Mailbox[*ClusterJoinNode]
				NotificationMailboxes *notificationMailboxes

				Delegate   *delegate
				Memberlist *memberlist.Memberlist
				Transport  *httpTransport
			}) Manager {
				worker := actor.NewWorker(func(ctx actor.Context) actor.WorkerStatus {
					select {
					case <-ctx.Done():
						return actor.WorkerEnd

					case msg, ok := <-d.NotificationMailboxes.publisher.ReceiveC():
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

					case msg, ok := <-d.NotificationMailboxes.consumer.ReceiveC():
						if !ok {
							return actor.WorkerEnd
						}

						if err := msg.Handle(ctx, d.Delegate); err != nil {
							d.Delegate.logger.ErrorContext(
								ctx,
								"failed to handle notification",
								slog.String("error", err.Error()),
							)
						}

						return actor.WorkerContinue

					case joinNode, ok := <-d.JoinM.ReceiveC():
						if !ok {
							return actor.WorkerEnd
						}

						d.Delegate.endpoints.Set(joinNode.JoinNodeID, joinNode.JoinNodeEndpoint)
						_, err := d.Memberlist.Join([]string{joinNode.Address()})
						if err != nil {
							d.Delegate.logger.ErrorContext(
								ctx,
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
					Actor:                  actor.New(worker),
					handler:                d.Transport,
					notificationPublisherM: d.NotificationMailboxes.publisher,
				}
			},
		),
		fxproviders.ProvideActor[*broadcaster](
			func(d struct {
				fx.In

				NotificationMailboxes *notificationMailboxes
				Notifier              changefeed.Notifier
			}) *broadcaster {
				return &broadcaster{
					Actor: actor.New(&broadcasterWorker{
						localNodeID:            opts.NodeID,
						notifier:               d.Notifier,
						notificationPublisherM: d.NotificationMailboxes.publisher,
					}),
				}
			},
		),
		fx.Invoke(func(_ struct {
			fx.In

			S *syncActor
			C *clusterFormationController
			B *broadcaster
		}) {
			// No-op, just to force the creation of all components
		}),
	)
}

func (m *managerImpl) HttpHandler() http.Handler {
	return m.handler
}
