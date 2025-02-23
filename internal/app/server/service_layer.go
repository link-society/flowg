package server

import (
	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"
)

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
) proctree.Process {
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

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		httpServer,
		mgmtServer,
		syslogServer,
	)
}
