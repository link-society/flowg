package server

import (
	"fmt"
	"strings"

	"crypto/tls"

	"go.uber.org/fx"

	"link-society.com/flowg/api/auth"

	"link-society.com/flowg/internal/app/config"
	"link-society.com/flowg/internal/app/featureflags"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/secret"

	badgerauth "link-society.com/flowg/internal/storage/backends/badger/concrete/auth"
	badgerconfig "link-society.com/flowg/internal/storage/backends/badger/concrete/config"
	badgerlog "link-society.com/flowg/internal/storage/backends/badger/concrete/log"
	fdbauth "link-society.com/flowg/internal/storage/backends/foundation/concrete/auth"
	fdbconfig "link-society.com/flowg/internal/storage/backends/foundation/concrete/config"
	fdblog "link-society.com/flowg/internal/storage/backends/foundation/concrete/log"

	"link-society.com/flowg/internal/services/http"
	"link-society.com/flowg/internal/services/mgmt"
	"link-society.com/flowg/internal/services/syslog"

	"link-society.com/flowg/web"
)

// configModule supplies the configuration to the graph, along with the
// providers converting it into the options of each component. It is the only
// place knowing both the configuration schema and the components' options.
func configModule(cfg *config.RootConfig) fx.Option {
	return fx.Module(
		"app.server.config",
		fx.Supply(cfg),
		fx.Provide(
			provideAuthKeys,
			provideFeatureFlags,
			provideHttpOptions,
			provideMgmtOptions,
			provideSyslogOptions,
			provideBootstrapOptions,
			provideWebOptions,
		),
	)
}

// storageModules returns the fx modules providing the auth, config and log
// storages of the backend selected by the configuration.
func storageModules(cfg *config.RootConfig) fx.Option {
	backend := cfg.Storage.Backend

	switch backend.Backend {
	case config.StorageBackendBadgerDb:
		return badgerDbModules()

	case config.StorageBackendFoundationDb:
		return foundationDbModules()

	default:
		return fx.Error(fmt.Errorf("unsupported storage backend: %q", backend.Backend))
	}
}

// badgerDbModules returns the fx modules providing the three storages backed by
// BadgerDB, along with the providers of their options.
func badgerDbModules() fx.Option {
	return fx.Options(
		fx.Provide(
			provideBadgerAuthOptions,
			provideBadgerConfigOptions,
			provideBadgerLogOptions,
		),
		badgerauth.Module(),
		badgerconfig.Module(),
		badgerlog.Module(),
	)
}

// provideBadgerAuthOptions configures the BadgerDB auth storage, stored in its
// own directory.
func provideBadgerAuthOptions(cfg *config.RootConfig) badgerauth.Options {
	opts := badgerauth.DefaultOptions()
	opts.Directory = cfg.Storage.Backend.BadgerDB.AuthDir
	return opts
}

// provideBadgerConfigOptions configures the BadgerDB config storage, stored in
// its own directory.
func provideBadgerConfigOptions(cfg *config.RootConfig) badgerconfig.Options {
	opts := badgerconfig.DefaultOptions()
	opts.Directory = cfg.Storage.Backend.BadgerDB.ConfigDir
	return opts
}

// provideBadgerLogOptions configures the BadgerDB log storage, stored in its
// own directory.
func provideBadgerLogOptions(cfg *config.RootConfig, flags featureflags.Flags) badgerlog.Options {
	opts := badgerlog.DefaultOptions()
	opts.Directory = cfg.Storage.Backend.BadgerDB.LogDir
	opts.StreamConfigOverride = streamConfigOverride(flags)
	return opts
}

// foundationDbModules returns the fx modules providing the three storages
// backed by FoundationDB, along with the providers of their options. The three
// storages share the same cluster and key space.
func foundationDbModules() fx.Option {
	return fx.Options(
		fx.Provide(
			provideFoundationDbAuthOptions,
			provideFoundationDbConfigOptions,
			provideFoundationDbLogOptions,
		),
		fdbauth.Module(),
		fdbconfig.Module(),
		fdblog.Module(),
	)
}

// provideFoundationDbAuthOptions configures the FoundationDB auth storage.
func provideFoundationDbAuthOptions(cfg *config.RootConfig) fdbauth.Options {
	fdbCfg := cfg.Storage.Backend.FoundationDB

	opts := fdbauth.DefaultOptions()
	opts.ClusterFile = fdbCfg.ClusterFile
	opts.ConnectionString = fdbCfg.ConnectionString
	if fdbCfg.KeySpace != "" {
		opts.KeySpace = fdbCfg.KeySpace
	}

	return opts
}

// provideFoundationDbConfigOptions configures the FoundationDB config storage.
func provideFoundationDbConfigOptions(cfg *config.RootConfig) fdbconfig.Options {
	fdbCfg := cfg.Storage.Backend.FoundationDB

	opts := fdbconfig.DefaultOptions()
	opts.ClusterFile = fdbCfg.ClusterFile
	opts.ConnectionString = fdbCfg.ConnectionString
	if fdbCfg.KeySpace != "" {
		opts.KeySpace = fdbCfg.KeySpace
	}

	return opts
}

// provideFoundationDbLogOptions configures the FoundationDB log storage.
func provideFoundationDbLogOptions(cfg *config.RootConfig, flags featureflags.Flags) fdblog.Options {
	fdbCfg := cfg.Storage.Backend.FoundationDB

	opts := fdblog.DefaultOptions()
	opts.ClusterFile = fdbCfg.ClusterFile
	opts.ConnectionString = fdbCfg.ConnectionString
	if fdbCfg.KeySpace != "" {
		opts.KeySpace = fdbCfg.KeySpace
	}
	opts.StreamConfigOverride = streamConfigOverride(flags)

	return opts
}

// streamConfigOverride returns the configuration enforced on every stream, or
// nil to use the stored ones. In demo mode, streams keep their logs for 15
// minutes and up to 10MB, without any index, so a public instance stays small.
func streamConfigOverride(flags featureflags.Flags) *models.StreamConfig {
	if !flags.DemoMode {
		return nil
	}

	return &models.StreamConfig{
		RetentionTime: 15 * 60, // 15 minutes
		RetentionSize: 10,
		IndexedFields: []string{},
	}
}

// provideAuthKeys provides the keys signing the JWTs issued on login, from the
// configured secret key. When no secret key is configured, a random one is
// generated, so issued tokens do not survive a restart.
func provideAuthKeys(cfg *config.RootConfig) (auth.Keys, error) {
	secretKey := cfg.Security.SecretKey
	if secretKey == "" {
		var err error
		secretKey, err = secret.NewSecret("jwt", 32)
		if err != nil {
			return auth.Keys{}, fmt.Errorf("failed to generate JWT secret key: %w", err)
		}
	}

	return auth.Keys{
		SigningKey: []byte(secretKey),
		TokenTTL:   auth.DefaultTokenTTL,
	}, nil
}

// provideFeatureFlags extracts the feature toggles enabled by the
// configuration.
func provideFeatureFlags(cfg *config.RootConfig) featureflags.Flags {
	return featureflags.Flags{
		DemoMode: cfg.DemoMode,
	}
}

// provideHttpOptions converts the configuration of the HTTP service into its
// options, loading its TLS certificate when TLS is enabled. Verbose logging
// also dumps the response bodies of failed requests.
func provideHttpOptions(cfg *config.RootConfig) (http.Options, error) {
	opts := http.Options{
		BindAddress:     cfg.Services.Http.BindAddress,
		MountPath:       mountPath(cfg),
		DumpErrorBodies: cfg.Logging.Verbose,
	}

	if tlsCfg := cfg.Services.Http.Tls; tlsCfg != nil {
		tlsConfig, err := loadTlsConfig("HTTP", tlsCfg.Cert, tlsCfg.Key)
		if err != nil {
			return http.Options{}, err
		}

		opts.TlsConfig = tlsConfig
	}

	return opts, nil
}

// provideMgmtOptions converts the configuration of the management service into
// its options, loading its TLS certificate when TLS is enabled.
func provideMgmtOptions(cfg *config.RootConfig) (mgmt.Options, error) {
	opts := mgmt.Options{
		BindAddress: cfg.Services.Management.BindAddress,
	}

	if tlsCfg := cfg.Services.Management.Tls; tlsCfg != nil {
		tlsConfig, err := loadTlsConfig("Management", tlsCfg.Cert, tlsCfg.Key)
		if err != nil {
			return mgmt.Options{}, err
		}

		opts.TlsConfig = tlsConfig
	}

	return opts, nil
}

// provideSyslogOptions converts the configuration of the syslog service into
// its options, loading its TLS certificate when TLS is enabled. Client
// certificates are verified when given, and required when TLS authentication is
// enabled.
func provideSyslogOptions(cfg *config.RootConfig) (syslog.Options, error) {
	opts := syslog.Options{
		TcpMode:     cfg.Services.Syslog.Protocol == "tcp",
		BindAddress: cfg.Services.Syslog.Bind,
	}

	if tlsCfg := cfg.Services.Syslog.Tls; tlsCfg != nil {
		tlsConfig, err := loadTlsConfig("Syslog", tlsCfg.Cert, tlsCfg.Key)
		if err != nil {
			return syslog.Options{}, err
		}

		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
		if tlsCfg.AuthEnabled {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		opts.TlsConfig = tlsConfig
	}

	return opts, nil
}

// provideWebOptions converts the configuration of the HTTP service into the
// options of the web UI it serves.
func provideWebOptions(cfg *config.RootConfig) web.Options {
	return web.Options{
		MountPath: mountPath(cfg),
	}
}

// provideBootstrapOptions extracts the settings applied when seeding the
// storage: the initial syslog origins, the initial/reset credentials, and
// whether to seed the demo role and user.
func provideBootstrapOptions(cfg *config.RootConfig, flags featureflags.Flags) bootstrapOptions {
	return bootstrapOptions{
		InitialSyslogAllowedOrigins: cfg.Services.Syslog.InitialAllowedOrigins,

		InitialUser:     cfg.Storage.Seed.Auth.InitialUser,
		InitialPassword: cfg.Storage.Seed.Auth.InitialPassword,

		DemoMode: flags.DemoMode,

		ResetUser:     cfg.Storage.Seed.Auth.ResetUser,
		ResetPassword: cfg.Storage.Seed.Auth.ResetPassword,
	}
}

// mountPath returns the path prefix the HTTP service is mounted under, without
// trailing slash ("" when mounted at the root).
func mountPath(cfg *config.RootConfig) string {
	return strings.TrimSuffix(cfg.Services.Http.MountPath, "/")
}

// loadTlsConfig loads the certificate of the named service into a TLS
// configuration requiring TLS 1.3.
func loadTlsConfig(service string, certFile string, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s TLS certificate: %w", service, err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}, nil
}
