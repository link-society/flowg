package http

import (
	"log/slog"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/utils/sync"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type Server struct {
	worker *worker
	actor  actor.Actor
}

func NewServer(
	bindAddress string,
	authStorage *auth.Storage,
	configStorage *config.Storage,
	logStorage *log.Storage,
	logNotifier *lognotify.LogNotifier,
	pipelineRunner *pipelines.Runner,
) *Server {
	worker := &worker{
		logger: slog.Default().With(slog.String("channel", "http")),

		authStorage:   authStorage,
		configStorage: configStorage,
		logStorage:    logStorage,

		logNotifier:    logNotifier,
		pipelineRunner: pipelineRunner,

		state:     &workerStarting{bindAddress: bindAddress},
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
