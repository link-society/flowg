package cluster

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

type LocalEndpointResolverCallback func() (*url.URL, error)

type ClusterFormationStrategy interface {
	Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error)
	Leave(ctx context.Context) error
}

type clusterFormationController struct {
	logger   *slog.Logger
	joinM    actor.MailboxSender[*ClusterJoinNode]
	resolver LocalEndpointResolverCallback
	strategy ClusterFormationStrategy
}

var _ actor.Worker = (*clusterFormationController)(nil)

func (c *clusterFormationController) DoWork(ctx actor.Context) actor.WorkerStatus {
	c.joiner(ctx)

	select {
	case <-ctx.Done():
		c.strategy.Leave(ctx)
		return actor.WorkerEnd

	case <-time.After(5 * time.Second):
		c.joiner(ctx)
		return actor.WorkerContinue
	}
}

func (c *clusterFormationController) joiner(ctx context.Context) {
	joinNodes, err := c.strategy.Join(ctx, c.resolver)
	if err != nil {
		c.logger.WarnContext(
			ctx,
			"failed to join cluster",
			slog.String("error", err.Error()),
		)
	} else {
		for _, joinNode := range joinNodes {
			c.joinM.Send(ctx, joinNode)
		}
	}
}
