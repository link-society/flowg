package server

import (
	"log/slog"

	"os"
	"os/signal"
	"syscall"

	"link-society.com/flowg/internal/server/service"
	"link-society.com/flowg/internal/server/storage"
)

type Options struct {
	HttpBindAddress   string
	SyslogBindAddress string

	ConfigStorageDir string
	AuthStorageDir   string
	LogStorageDir    string
}

func Run(opts Options, errC chan<- struct{}) {
	defer close(errC)

	storageA := storage.New(
		opts.AuthStorageDir,
		opts.LogStorageDir,
		opts.ConfigStorageDir,
	)
	storageA.Start()
	defer func() {
		storageA.Stop()

		for err := range storageA.StopErrC() {
			errC <- err
		}
	}()

	erred := false
	for err := range storageA.StartErrC() {
		errC <- err
		erred = true
	}

	if erred {
		return
	}

	serviceA := service.New(
		opts.HttpBindAddress,
		opts.SyslogBindAddress,
		storageA,
	)
	serviceA.Start()
	defer func() {
		serviceA.Stop()

		for err := range serviceA.StopErrC() {
			errC <- err
		}
	}()

	erred = false
	for err := range serviceA.StartErrC() {
		errC <- err
		erred = true
	}

	if erred {
		return
	}

	slog.Info(
		"server ready",
		"channel", "server",
	)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
	<-sigC
}
