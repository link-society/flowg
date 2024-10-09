package syslog

import (
	"log/slog"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/storage/config"
)

type Server struct {
	worker *worker
	actor  actor.Actor
}

func NewServer(
	isTCP bool,
	bindAddress string,
	tlsConfig *tls.Config,
	allowOrigins []string,

	configStorage *config.Storage,
	pipelineRunner *pipelines.Runner,
) *Server {
	worker := &worker{
		logger: slog.Default().With(slog.String("channel", "syslog")),

		configStorage:  configStorage,
		pipelineRunner: pipelineRunner,

		state: &workerStarting{
			isTCP:        isTCP,
			bindAddress:  bindAddress,
			tlsConfig:    tlsConfig,
			allowOrigins: allowOrigins,
		},

		startCond: sync.NewCondValue[error](),
		stopCond:  sync.NewCondValue[error](),
	}

	actor := actor.New(worker)

	return &Server{worker, actor}
}

func (s *Server) Start() {
	s.actor.Start()
}

func (s *Server) WaitStarted() error {
	return s.worker.startCond.Wait()
}

func (s *Server) Stop() {
	s.actor.Stop()
}

func (s *Server) WaitStopped() error {
	return s.worker.stopCond.Wait()
}
