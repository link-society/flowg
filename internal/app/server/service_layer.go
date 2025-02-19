package server

import (
	"errors"

	"crypto/tls"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
)

type serviceLayer struct {
	httpServer   *http.Server
	mgmtServer   *mgmt.Server
	syslogServer *syslog.Server

	actor actor.Actor
}

func newServiceLayer(
	httpBindAddress string,
	httpTlsConfig *tls.Config,

	mgmtBindAddress string,
	mgmtTlsConfig *tls.Config,

	syslogTCP bool,
	syslogBindAddress string,
	syslogTlsConfig *tls.Config,
	syslogAllowOrigins []string,

	storageLayer *storageLayer,
	engineLayer *engineLayer,
) *serviceLayer {
	httpServer := http.NewServer(
		httpBindAddress,
		httpTlsConfig,
		storageLayer.authStorage,
		storageLayer.configStorage,
		storageLayer.logStorage,
		engineLayer.logNotifier,
		engineLayer.pipelineRunner,
	)

	mgmtServer := mgmt.NewServer(
		mgmtBindAddress,
		mgmtTlsConfig,
	)

	syslogServer := syslog.NewServer(
		syslogTCP,
		syslogBindAddress,
		syslogTlsConfig,
		syslogAllowOrigins,

		storageLayer.configStorage,
		engineLayer.pipelineRunner,
	)

	rootA := actor.Combine(httpServer, mgmtServer, syslogServer).
		WithOptions(actor.OptStopTogether()).
		Build()

	return &serviceLayer{
		httpServer:   httpServer,
		mgmtServer:   mgmtServer,
		syslogServer: syslogServer,

		actor: rootA,
	}
}

func (s *serviceLayer) Start() {
	s.actor.Start()
}

func (s *serviceLayer) WaitStarted() error {
	errs := []error{}

	if err := s.httpServer.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if err := s.mgmtServer.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if err := s.syslogServer.WaitStarted(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.New("some services failed to start")
	}

	return nil
}

func (s *serviceLayer) Stop() {
	s.actor.Stop()
}

func (s *serviceLayer) WaitStopped() error {
	errs := []error{}

	if err := s.syslogServer.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if err := s.mgmtServer.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if err := s.httpServer.WaitStopped(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.New("some services failed to stop")
	}

	return nil
}
