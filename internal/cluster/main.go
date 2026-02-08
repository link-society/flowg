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

	"github.com/hashicorp/memberlist"

	clusterstate "link-society.com/flowg/internal/storage/cluster-state"
)

type Manager interface {
	HttpHandler() http.Handler
}

type ManagerOptions struct {
	NodeID string
	Cookie string

	ClusterFormationStrategy ClusterFormationStrategy
	ClusterStateDir          string
}

type managerImpl struct {
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
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[notification] {
			mbox := actor.NewMailbox[notification]()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					mbox.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					mbox.Stop()
					return nil
				},
			})

			return mbox
		}),
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[net.Conn] {
			mbox := actor.NewMailbox[net.Conn]()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					mbox.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					mbox.Stop()
					return nil
				},
			})

			return mbox
		}),
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[*memberlist.Packet] {
			mbox := actor.NewMailbox[*memberlist.Packet]()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					mbox.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					mbox.Stop()
					return nil
				},
			})

			return mbox
		}),
		fx.Provide(func(lc fx.Lifecycle) actor.Mailbox[*ClusterJoinNode] {
			mbox := actor.NewMailbox[*ClusterJoinNode]()

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					mbox.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					mbox.Stop()
					return nil
				},
			})

			return mbox
		}),
		fx.Provide(fx.Annotate(
			func(
				lc fx.Lifecycle,
				joinM actor.Mailbox[*ClusterJoinNode],
				listener *Listener,
			) actor.Actor {
				a := actor.New(&clusterFormationController{
					logger:   slog.Default().With(slog.String("channel", "cluster.formation")),
					joinM:    joinM,
					resolver: func() (*url.URL, error) { return listener.ResolveLocalEndpoint() },
					strategy: opts.ClusterFormationStrategy,
				})

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						a.Start()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						a.Stop()
						return nil
					},
				})

				return a
			},
			fx.ResultTags(`name:"cluster.manager.formation"`),
		)),
		fx.Provide(func(listener *Listener) (*delegate, error) {
			var err error

			localEndpoint, err := listener.ResolveLocalEndpoint()
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
			}, nil
		}),
		fx.Provide(func(
			delegate *delegate,
			connM actor.Mailbox[net.Conn],
			packetM actor.Mailbox[*memberlist.Packet],
		) *httpTransport {
			return &httpTransport{
				delegate: delegate,
				cookie:   opts.Cookie,

				connM:   connM,
				packetM: packetM,
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
		fx.Provide(func(d struct {
			fx.In

			LC fx.Lifecycle

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
			a := actor.New(worker)

			d.LC.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					a.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					a.Stop()
					return nil
				},
			})

			return &managerImpl{handler: d.Transport, notifyM: d.NotifyM}
		}),
		fx.Invoke(func(_ struct {
			fx.In
			C actor.Actor `name:"cluster.manager.formation"`
		}) {
			// No-op, just to force the creation of all components
		}),
	)
}

func (m *managerImpl) HttpHandler() http.Handler {
	return m.handler
}
