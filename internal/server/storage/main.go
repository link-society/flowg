package storage

import (
	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/data/lognotify"
)

type Actor struct {
	AuthA        *authA
	LogA         *logA
	ConfigA      *configA
	LogNotifierA *lognotify.LogNotifier

	rootA     actor.Actor
	startErrC chan struct{}
	stopErrC  chan struct{}
}

func New(authDir, logDir, configDir string) *Actor {
	var (
		authA        = newAuthA(authDir)
		logA         = newLogA(logDir)
		configA      = newConfigA(configDir)
		logNotifierA = lognotify.NewLogNotifier()
	)

	rootA := actor.Combine(authA, logA, configA, logNotifierA).Build()

	return &Actor{
		AuthA:        authA,
		LogA:         logA,
		ConfigA:      configA,
		LogNotifierA: logNotifierA,

		rootA:     rootA,
		startErrC: make(chan struct{}, 3),
		stopErrC:  make(chan struct{}, 3),
	}
}

func (a *Actor) Start() {
	a.rootA.Start()

	if err, ok := <-a.AuthA.StartErrC(); ok {
		a.startErrC <- err
	}

	if err, ok := <-a.LogA.StartErrC(); ok {
		a.startErrC <- err
	}

	if err, ok := <-a.ConfigA.StartErrC(); ok {
		a.startErrC <- err
	}

	close(a.startErrC)
}

func (a *Actor) Stop() {
	a.rootA.Stop()

	if err, ok := <-a.AuthA.StopErrC(); ok {
		a.stopErrC <- err
	}

	if err, ok := <-a.LogA.StopErrC(); ok {
		a.stopErrC <- err
	}

	if err, ok := <-a.ConfigA.StopErrC(); ok {
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
