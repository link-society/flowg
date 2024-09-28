package syslog

import (
	"log/slog"

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
	bindAddress string,
	configStorage *config.Storage,
	pipelineRunner *pipelines.Runner,
) *Server {
	worker := &worker{
		logger: slog.Default().With(slog.String("channel", "syslog")),

		configStorage:  configStorage,
		pipelineRunner: pipelineRunner,

		state: &workerStarting{bindAddress: bindAddress},

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
