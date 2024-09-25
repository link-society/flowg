package service

import (
	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/integrations/syslog"

	"link-society.com/flowg/internal/server/storage"
)

type Actor struct {
	httpA   *httpA
	syslogA actor.Actor

	rootA     actor.Actor
	startErrC chan struct{}
	stopErrC  chan struct{}
}

func New(
	httpBindAddress string,
	syslogBindAddress string,
	storageA *storage.Actor,
) *Actor {
	var (
		httpA = newHttpA(
			httpBindAddress,
			storageA.AuthA.Database,
			storageA.LogA.Storage,
			storageA.ConfigA.Storage,
			storageA.LogNotifierA,
		)
		syslogA = syslog.NewServer(
			syslogBindAddress,
			storageA.ConfigA.Storage,
			storageA.LogA.Storage,
			storageA.LogNotifierA,
		)
	)

	rootA := actor.Combine(httpA, syslogA).Build()

	return &Actor{
		httpA:   httpA,
		syslogA: syslogA,

		rootA:     rootA,
		startErrC: make(chan struct{}, 1),
		stopErrC:  make(chan struct{}, 1),
	}
}

func (a *Actor) Start() {
	a.rootA.Start()

	if err, ok := <-a.httpA.StartErrC(); ok {
		a.startErrC <- err
	}

	close(a.startErrC)
}

func (a *Actor) Stop() {
	a.rootA.Stop()

	if err, ok := <-a.httpA.StopErrC(); ok {
		a.stopErrC <- err
	}

	close(a.stopErrC)
}

func (a *Actor) StartErrC() <-chan struct{} {
	return a.startErrC
}

func (a *Actor) StopErrC() <-chan struct{} {
	return a.stopErrC
}
