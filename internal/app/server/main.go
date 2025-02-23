package server

import (
	"log/slog"

	"crypto/tls"

	"link-society.com/flowg/internal/utils/proctree"
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

func NewServer(opts Options) proctree.Process {
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

	bootstrap := proctree.NewProcess(&bootstrapProcHandler{
		logger:       slog.Default().With("channel", "server"),
		storageLayer: storageLayer,
	})

	return proctree.NewProcessGroup(
		proctree.DefaultProcessGroupOptions(),
		storageLayer,
		engineLayer,
		serviceLayer,
		bootstrap,
	)
}
