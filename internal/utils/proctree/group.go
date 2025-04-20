package proctree

import (
	"context"
	"errors"

	"time"
)

type ProcessGroupOptions struct {
	InitTimeout time.Duration
	JoinTimeout time.Duration
}

func DefaultProcessGroupOptions() ProcessGroupOptions {
	return ProcessGroupOptions{
		InitTimeout: 1 * time.Minute, // Automatic cluster formation may not be possible in default 5 seconds timeout
		JoinTimeout: 5 * time.Second,
	}
}

func NewProcessGroup(opts ProcessGroupOptions, children ...Process) Process {
	return &group{
		opts:     opts,
		children: children,

		readyCond: newCondValue[error](),
		joinCond:  newCondValue[error](),
		shutdownC: make(chan struct{}, 1),
	}
}

type group struct {
	opts     ProcessGroupOptions
	children []Process

	readyCond *condValue[error]
	joinCond  *condValue[error]
	shutdownC chan struct{}
}

var _ Process = (*group)(nil)

func (g *group) Start() {
	go g.run()
}

func (g *group) Stop() {
	g.shutdownC <- struct{}{}
}

func (g *group) WaitReady(ctx context.Context) error {
	readyC := make(chan error, 1)

	go func() {
		err := g.readyCond.Wait()
		readyC <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-readyC:
		return err
	}
}

func (g *group) Join(ctx context.Context) error {
	doneC := make(chan error, 1)

	go func() {
		err := g.joinCond.Wait()
		doneC <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-doneC:
		return err
	}
}

func (g *group) run() {
	for i, child := range g.children {
		child.Start()

		ctx, cancel := context.WithTimeout(context.Background(), g.opts.InitTimeout)
		err := child.WaitReady(ctx)
		cancel()

		if err != nil {
			errs := g.stopChildren(i)
			errs = append([]error{err}, errs...)
			err = errors.Join(errs...)
			g.readyCond.Broadcast(err)
			g.joinCond.Broadcast(err)
			return
		}
	}

	joinC := make(chan struct{}, len(g.children))
	monitorCtx, monitorCancel := context.WithCancel(context.Background())

	for _, child := range g.children {
		go func(child Process) {
			child.Join(monitorCtx)
			joinC <- struct{}{}
		}(child)
	}

	g.readyCond.Broadcast(nil)

	select {
	case <-g.shutdownC:
	case <-joinC:
	}

	monitorCancel()
	errs := g.stopChildren(len(g.children) - 1)
	err := errors.Join(errs...)
	g.joinCond.Broadcast(err)
	close(g.shutdownC)
}

func (g *group) stopChildren(last int) []error {
	errs := make([]error, 0, len(g.children)-last)

	for i := last; i >= 0; i-- {
		child := g.children[i]
		child.Stop()

		ctx, cancel := context.WithTimeout(context.Background(), g.opts.JoinTimeout)
		err := child.Join(ctx)
		cancel()

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
