package cmd

import (
	"fmt"
	"strings"

	"crypto/tls"

	"link-society.com/flowg/internal/app/config"
	"link-society.com/flowg/internal/app/server"
)

// Converts the parsed RootConfig into a server.Options struct, loading TLS
// certificates for the services that have TLS enabled.
func asServerOptions(cfg *config.RootConfig) (server.Options, error) {
	var (
		httpTlsConfig   *tls.Config
		mgmtTlsConfig   *tls.Config
		syslogTlsConfig *tls.Config
	)

	if cfg.Services.Http.Tls != nil {
		cert, err := tls.LoadX509KeyPair(cfg.Services.Http.Tls.Cert, cfg.Services.Http.Tls.Key)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load HTTP TLS certificate: %w", err)
		}

		httpTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}
	}

	if cfg.Services.Management.Tls != nil {
		cert, err := tls.LoadX509KeyPair(cfg.Services.Management.Tls.Cert, cfg.Services.Management.Tls.Key)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load Management TLS certificate: %w", err)
		}

		mgmtTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}
	}

	if cfg.Services.Syslog.Tls != nil {
		cert, err := tls.LoadX509KeyPair(cfg.Services.Syslog.Tls.Cert, cfg.Services.Syslog.Tls.Key)
		if err != nil {
			return server.Options{}, fmt.Errorf("failed to load Syslog TLS certificate: %w", err)
		}

		clientAuth := tls.VerifyClientCertIfGiven
		if cfg.Services.Syslog.Tls.AuthEnabled {
			clientAuth = tls.RequireAndVerifyClientCert
		}

		syslogTlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   clientAuth,
			MinVersion:   tls.VersionTLS13,
		}
	}

	var storageOptions server.StorageOptions
	switch cfg.Storage.Backend.Backend {
	case config.StorageBackendBadgerDb:
		storageOptions = &server.BadgerDbStorageOptions{
			AuthDir:   cfg.Storage.Backend.BadgerDB.AuthDir,
			LogDir:    cfg.Storage.Backend.BadgerDB.LogDir,
			ConfigDir: cfg.Storage.Backend.BadgerDB.ConfigDir,
		}

	case config.StorageBackendFoundationDb:
		storageOptions = &server.FoundationDbStorageOptions{
			ClusterFile:      cfg.Storage.Backend.FoundationDB.ClusterFile,
			ConnectionString: cfg.Storage.Backend.FoundationDB.ConnectionString,
			KeySpace:         cfg.Storage.Backend.FoundationDB.KeySpace,
		}
	}

	opts := server.Options{
		HttpBindAddress: cfg.Services.Http.BindAddress,
		HttpMountPath:   strings.TrimSuffix(cfg.Services.Http.MountPath, "/"),
		HttpTlsConfig:   httpTlsConfig,

		MgmtBindAddress: cfg.Services.Management.BindAddress,
		MgmtTlsConfig:   mgmtTlsConfig,

		SyslogTcpMode:               cfg.Services.Syslog.Protocol == "tcp",
		SyslogBindAddress:           cfg.Services.Syslog.Bind,
		SyslogTlsConfig:             syslogTlsConfig,
		SyslogInitialAllowedOrigins: cfg.Services.Syslog.InitialAllowedOrigins,

		Storage: storageOptions,

		AuthInitialUser:     cfg.Storage.Seed.Auth.InitialUser,
		AuthInitialPassword: cfg.Storage.Seed.Auth.InitialPassword,

		AuthResetUser:     cfg.Storage.Seed.Auth.ResetUser,
		AuthResetPassword: cfg.Storage.Seed.Auth.ResetPassword,
	}

	return opts, nil
}
