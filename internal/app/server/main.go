package server

import (
	"log/slog"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"
)

type Options struct {
	HttpBindAddress string
	HttpTlsConfig   *tls.Config

	MgmtBindAddress string
	MgmtTlsConfig   *tls.Config

	SyslogTCP          bool
	SyslogBindAddress  string
	SyslogTlsConfig    *tls.Config
	SyslogAllowOrigins []string

	AuthStorageDir   string
	ConfigStorageDir string
	LogStorageDir    string
}

type Server struct {
	actor actor.Actor
	doneC chan bool
}

func NewServer(opts Options) *Server {
	storageLayer := newStorageLayer(
		opts.AuthStorageDir,
		opts.ConfigStorageDir,
		opts.LogStorageDir,
	)
	engineLayer := newEngineLayer(
		storageLayer,
	)
	serviceLayer := newServiceLayer(
		opts.HttpBindAddress,
		opts.HttpTlsConfig,

		opts.MgmtBindAddress,
		opts.MgmtTlsConfig,

		opts.SyslogTCP,
		opts.SyslogBindAddress,
		opts.SyslogTlsConfig,
		opts.SyslogAllowOrigins,

		storageLayer,
		engineLayer,
	)

	doneC := make(chan bool, 1)

	worker := &worker{
		state:  &workerStartingStorageLayer{},
		logger: slog.Default().With("channel", "server"),

		storageLayer: storageLayer,
		engineLayer:  engineLayer,
		serviceLayer: serviceLayer,
	}

	rootA := actor.New(
		worker,
		actor.OptOnStop(func() {
			doneC <- worker.failure
			close(doneC)
		}),
	)

	return &Server{
		actor: rootA,
		doneC: doneC,
	}
}

func (s *Server) Start() {
	s.actor.Start()
}

func (s *Server) Stop() {
	s.actor.Stop()
}

func (s *Server) DoneC() <-chan bool {
	return s.doneC
}
