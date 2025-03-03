package proctree

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
)

func NewActorProcess(parent actor.Actor) Process {
	readyCond := newCondValue[struct{}]()
	doneCond := newCondValue[error]()

	waiter := actor.New(
		actor.NewWorker(func(ctx actor.Context) actor.WorkerStatus {
			<-ctx.Done()
			return actor.WorkerEnd
		}),
	)

	rootA := actor.Combine(parent, waiter).
		WithOptions(
			actor.OptStopTogether(),
			actor.OptOnStartCombined(func(ctx actor.Context) {
				readyCond.Broadcast(struct{}{})
			}),
			actor.OptOnStopCombined(func() {
				doneCond.Broadcast(nil)
			}),
		).
		Build()

	return &actorProcess{
		Actor:     rootA,
		readyCond: readyCond,
		doneCond:  doneCond,
	}
}

type actorProcess struct {
	actor.Actor

	readyCond *condValue[struct{}]
	doneCond  *condValue[error]
}

func (p *actorProcess) WaitReady(ctx context.Context) error {
	readyC := make(chan struct{}, 1)

	go func() {
		p.readyCond.Wait()
		readyC <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-readyC:
		return nil
	}
}

func (p *actorProcess) Join(ctx context.Context) error {
	doneC := make(chan error, 1)

	go func() {
		err := p.doneCond.Wait()
		doneC <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-doneC:
		return err
	}
}
