package mgmt

import (
	"log/slog"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/sync"
)

type Server struct {
	worker *worker
	actor  actor.Actor
}

func NewServer(
	bindAddress string,
	tlsConfig *tls.Config,
) *Server {
	worker := &worker{
		logger: slog.Default().With(slog.String("channel", "mgmt")),

		state: &workerStarting{
			bindAddress: bindAddress,
			tlsConfig:   tlsConfig,
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
