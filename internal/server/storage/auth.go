package storage

import (
	"log/slog"

	"link-society.com/flowg/internal/app/bootstrap"

	"link-society.com/flowg/internal/data/auth"
)

type authA struct {
	Database *auth.Database
	dir      string

	startErrC chan struct{}
	stopErrC  chan struct{}
}

func newAuthA(dir string) *authA {
	return &authA{
		Database: auth.NewDatabase(auth.DefaultDatabaseOpts().WithDir(dir)),
		dir:      dir,

		startErrC: make(chan struct{}, 1),
		stopErrC:  make(chan struct{}, 1),
	}
}

func (a *authA) Start() {
	defer close(a.startErrC)

	err := a.Database.Open()
	if err != nil {
		slog.Error(
			"Failed to open auth database",
			"channel", "auth",
			"path", a.dir,
			"error", err,
		)
		a.startErrC <- struct{}{}
		return
	}

	if err := bootstrap.DefaultRolesAndUsers(a.Database); err != nil {
		slog.Error(
			"Failed to bootstrap default roles and users",
			"channel", "auth",
			"error", err,
		)
		a.startErrC <- struct{}{}
		return
	}
}

func (a *authA) Stop() {
	defer close(a.stopErrC)

	if a.Database == nil {
		return
	}

	err := a.Database.Close()
	if err != nil {
		slog.Error(
			"Failed to close auth database",
			"channel", "auth",
			"path", a.dir,
			"error", err,
		)
		a.stopErrC <- struct{}{}
		return
	}
}

func (a *authA) StartErrC() <-chan struct{} {
	return a.startErrC
}

func (a *authA) StopErrC() <-chan struct{} {
	return a.stopErrC
}
