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

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
	"link-society.com/flowg/internal/utils/kvstore"
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
}

type deps struct {
	fx.In

	ClusterStateStorage kvstore.Storage `name:"cluster.state"`
	AuthStorage         auth.Storage
	ConfigStorage       config.Storage
	LogStorage          log.Storage
}

var _ Manager = (*managerImpl)(nil)

func NewManager(opts ManagerOptions) fx.Option {
	kvOpts := kvstore.DefaultOptions()
	kvOpts.LogChannel = "cluster.state"
	kvOpts.Directory = opts.ClusterStateDir

	return fx.Module(
		"cluster.manager",
		kvstore.NewStorage(kvOpts),
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
		fx.Provide(func(deps deps, listener *Listener) (*delegate, error) {
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
				endpoints:     make(map[string]*url.URL),

				clusterStateStorage: deps.ClusterStateStorage,

				syncPool: &syncPool{
					logger: slog.Default().With(
						slog.String("channel", "cluster.replication"),
						slog.String("cluster.local.node", opts.NodeID),
						slog.String("cluster.local.endpoint", localEndpoint.String()),
					),

					nodeID: opts.NodeID,
					cookie: opts.Cookie,

					authStorage:   deps.AuthStorage,
					configStorage: deps.ConfigStorage,
					logStorage:    deps.LogStorage,

					workers: make(map[string]*syncActor),
				},
			}, nil
		}),
		fx.Provide(func(
			delegate *delegate,
			connM actor.Mailbox[net.Conn],
			packetM actor.Mailbox[*memberlist.Packet],
			deps deps,
		) *httpTransport {
			return &httpTransport{
				delegate: delegate,
				cookie:   opts.Cookie,

				connM:   connM,
				packetM: packetM,

				authStorage:   deps.AuthStorage,
				configStorage: deps.ConfigStorage,
				logStorage:    deps.LogStorage,
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
					delegate.syncPool.RemoveAll()

					if err := mlist.Leave(5 * time.Second); err != nil {
						return err
					}

					return mlist.Shutdown()
				},
			})

			return mlist, nil
		}),
		fx.Provide(func(
			lc fx.Lifecycle,
			joinM actor.Mailbox[*ClusterJoinNode],
			delegate *delegate,
			mlist *memberlist.Memberlist,
			transport *httpTransport,
		) Manager {
			worker := actor.NewWorker(func(ctx actor.Context) actor.WorkerStatus {
				select {
				case <-ctx.Done():
					return actor.WorkerEnd

				case joinNode, ok := <-joinM.ReceiveC():
					if !ok {
						return actor.WorkerEnd
					}

					delegate.endpoints[joinNode.JoinNodeID] = joinNode.JoinNodeEndpoint
					_, err := mlist.Join([]string{joinNode.Address()})
					if err != nil {
						delegate.logger.Error(
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

			return &managerImpl{handler: transport}
		}),
	)
}

func (m *managerImpl) HttpHandler() http.Handler {
	return m.handler
}
