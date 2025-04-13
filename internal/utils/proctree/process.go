package proctree

import (
	"context"

	"github.com/vladopajic/go-actor/actor"
)

type Process interface {
	actor.Actor

	WaitReady(ctx context.Context) error
	Join(ctx context.Context) error
}

type ProcessHandler interface {
	Init(ctx actor.Context) ProcessResult
	DoWork(ctx actor.Context) ProcessResult
	Terminate(ctx actor.Context, err error) error
}

func NewProcess(handler ProcessHandler) Process {
	worker := &procWorker{
		handler:   handler,
		state:     &procWorkerInit{},
		readyCond: newCondValue[error](),
		joinCond:  newCondValue[error](),
	}

	return &proc{
		Actor:  actor.New(worker),
		worker: worker,
	}
}

type proc struct {
	actor.Actor
	worker *procWorker
}

var _ Process = (*proc)(nil)

func (p *proc) WaitReady(ctx context.Context) error {
	readyC := make(chan error, 1)

	go func() {
		err := p.worker.readyCond.Wait()
		readyC <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-readyC:
		return err
	}
}

func (p *proc) Join(ctx context.Context) error {
	doneC := make(chan error, 1)

	go func() {
		err := p.worker.joinCond.Wait()
		doneC <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case err := <-doneC:
		return err
	}
}

type procWorker struct {
	handler   ProcessHandler
	state     procWorkerState
	readyCond *condValue[error]
	joinCond  *condValue[error]
}

func (worker *procWorker) DoWork(ctx actor.Context) actor.WorkerStatus {
	worker.state = worker.state.DoWork(ctx, worker)

	if worker.state == nil {
		return actor.WorkerEnd
	}

	return actor.WorkerContinue
}

type procWorkerState interface {
	DoWork(ctx actor.Context, worker *procWorker) procWorkerState
}

type procWorkerInit struct{}

func (s *procWorkerInit) DoWork(ctx actor.Context, worker *procWorker) procWorkerState {
	if result := worker.handler.Init(ctx); result.Done() {
		worker.readyCond.Broadcast(result.Error())
		return &procWorkerTerminate{err: result.Error()}
	}

	worker.readyCond.Broadcast(nil)
	return &procWorkerRunning{}
}

type procWorkerRunning struct{}

func (s *procWorkerRunning) DoWork(ctx actor.Context, worker *procWorker) procWorkerState {
	if result := worker.handler.DoWork(ctx); result.Done() {
		return &procWorkerTerminate{err: result.Error()}
	}

	return s
}

type procWorkerTerminate struct {
	err error
}

func (s *procWorkerTerminate) DoWork(ctx actor.Context, worker *procWorker) procWorkerState {
	err := worker.handler.Terminate(ctx, s.err)
	worker.joinCond.Broadcast(err)
	return nil
}
