package cmd

import (
	"errors"

	"fmt"
	"strings"

	"crypto/tls"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"

	"link-society.com/flowg/internal/app/server"
)

// Represents the configuration file structure
type RootConfig struct {
	// DemoMode enables the demo role and user account which is not suitable for
	// production use. It is intended for demonstration and testing purposes only.
	DemoMode bool `hcl:"demo,optional"`

	Logging  *LoggingConfig  `hcl:"logging,block"`
	Services *ServicesConfig `hcl:"services,block"`
	Storage  *StorageConfig  `hcl:"storage,block"`
}

// Define logging settings
type LoggingConfig struct {
	Level   string `hcl:"level,optional"`
	Verbose bool   `hcl:"verbose,optional"`
}

// Group together the configuration for the services provided by the server.
// Each service needs to bind to a specific address and port, and may have TLS
// support.
type ServicesConfig struct {
	Http       *HttpConfig       `hcl:"http,block"`
	Management *ManagementConfig `hcl:"management,block"`
	Syslog     *SyslogConfig     `hcl:"syslog,block"`
}

// Configuration for the "HTTP" service, which provides the REST API and web UI.
type HttpConfig struct {
	BindAddress string         `hcl:"bind,optional"`
	MountPath   string         `hcl:"mount,optional"`
	Tls         *HttpTlsConfig `hcl:"tls,block"`
}

// Configuration for the "Management" service, which provides the health check
// and metrics endpoints.
type ManagementConfig struct {
	BindAddress string         `hcl:"bind,optional"`
	Tls         *HttpTlsConfig `hcl:"tls,block"`
}

// Configuration for TLS settings for HTTP services.
type HttpTlsConfig struct {
	Cert string `hcl:"cert"`
	Key  string `hcl:"key"`
}

// Configuration for the "Syslog" service, which provides a syslog server for
// receiving log messages from other services. It supports both TCP and UDP
// protocols, and can optionally use TLS for secure communication.
type SyslogConfig struct {
	Bind                  string           `hcl:"bind,optional"`
	Protocol              string           `hcl:"protocol,optional"`
	InitialAllowedOrigins []string         `hcl:"initial_allowed_origins,optional"`
	Tls                   *SyslogTlsConfig `hcl:"tls,block"`
}

// Configuration for TLS settings for the Syslog service.
type SyslogTlsConfig struct {
	Cert        string `hcl:"cert"`
	Key         string `hcl:"key"`
	AuthEnabled bool   `hcl:"auth,optional"`
}

// StorageConfig groups together settings related to storage.
type StorageConfig struct {
	Backend *StorageBackendConfig `hcl:"backend,block"`
	Seed    *StorageSeedConfig    `hcl:"seed,block"`
}

// Configuration for the storage backend used by the server. It supports
// multiple backends, but only one can be configured at a time.
type StorageBackendConfig struct {
	Backend string   `hcl:"backend,label"`
	Body    hcl.Body `hcl:",remain"`

	BadgerDB *StorageBackendBadgerDbConfig
}

// Constants defining the supported storage backends.
const (
	StorageBackendBadgerDb     = "badgerdb"
	StorageBackendFoundationDb = "foundationdb"
)

// Errors related to storage backend configuration.
var (
	ErrNoStorageBackend        = errors.New("no storage backend configured")
	ErrMultipleStorageBackends = errors.New("multiple storage backends configured")
)

// Configuration for the BadgerDB storage backend, which is an embedded
// key-value store.
type StorageBackendBadgerDbConfig struct {
	AuthDir   string `hcl:"auth_dir,optional"`
	LogDir    string `hcl:"logs_dir,optional"`
	ConfigDir string `hcl:"config_dir,optional"`
}

// Configuration for seeding the storage backend with initial data. This is
// useful for setting up default users, roles, and other necessary data during
// the first run of the server.
type StorageSeedConfig struct {
	Auth *StorageSeedAuthConfig `hcl:"auth,block"`
}

// Configuration for seeding the authentication storage with an initial user
// and password. This is useful for bootstrapping the server with a default
// admin account.
type StorageSeedAuthConfig struct {
	InitialUser     string `hcl:"initial_user"`
	InitialPassword string `hcl:"initial_password"`

	ResetUser     string `hcl:"reset_user,optional"`
	ResetPassword string `hcl:"reset_password,optional"`
}

// Returns a default configuration for the server, which is inherited from
// environment variables.
func DefaultConfig() *RootConfig {
	var defaultHttpTlsConfig *HttpTlsConfig
	if defaultHttpTlsEnabled {
		defaultHttpTlsConfig = &HttpTlsConfig{
			Cert: defaultHttpTlsCert,
			Key:  defaultHttpTlsCertKey,
		}
	}

	var defaultMgmtTlsConfig *HttpTlsConfig
	if defaultMgmtTlsEnabled {
		defaultMgmtTlsConfig = &HttpTlsConfig{
			Cert: defaultMgmtTlsCert,
			Key:  defaultMgmtTlsCertKey,
		}
	}

	var defaultSyslogTlsConfig *SyslogTlsConfig
	if defaultSyslogTlsEnabled {
		defaultSyslogTlsConfig = &SyslogTlsConfig{
			Cert:        defaultSyslogTlsCert,
			Key:         defaultSyslogTlsCertKey,
			AuthEnabled: defaultSyslogTlsAuthEnabled,
		}
	}

	storageBackendConfig := &StorageBackendConfig{
		Backend: defaultStorageBackend,
	}

	switch storageBackendConfig.Backend {
	case StorageBackendBadgerDb:
		storageBackendConfig.BadgerDB = DefaultStorageBackendBadgerDbConfig()

	case StorageBackendFoundationDb:
		panic("not implemented")
	}

	return &RootConfig{
		DemoMode: defaultDemoMode,
		Logging: &LoggingConfig{
			Level:   defaultLogLevel,
			Verbose: defaultVerbose,
		},
		Services: &ServicesConfig{
			Http: &HttpConfig{
				BindAddress: defaultHttpBindAddress,
				MountPath:   defaultHttpMountPath,
				Tls:         defaultHttpTlsConfig,
			},
			Management: &ManagementConfig{
				BindAddress: defaultMgmtBindAddress,
				Tls:         defaultMgmtTlsConfig,
			},
			Syslog: &SyslogConfig{
				Bind:     defaultSyslogBindAddr,
				Protocol: defaultSyslogProtocol,
				Tls:      defaultSyslogTlsConfig,
			},
		},
		Storage: &StorageConfig{
			Backend: storageBackendConfig,
			Seed: &StorageSeedConfig{
				Auth: &StorageSeedAuthConfig{
					InitialUser:     defaultAuthInitialUser,
					InitialPassword: defaultAuthInitialPassword,
					ResetUser:       defaultAuthResetUser,
					ResetPassword:   defaultAuthResetPassword,
				},
			},
		},
	}
}

// Returns a default configuration for the BadgerDB storage backend, which is
// inherited from environment variables.
func DefaultStorageBackendBadgerDbConfig() *StorageBackendBadgerDbConfig {
	return &StorageBackendBadgerDbConfig{
		AuthDir:   defaultBadgerAuthDir,
		LogDir:    defaultBadgerLogDir,
		ConfigDir: defaultBadgerConfigDir,
	}
}

// Validates the root configuration, ensuring that all required fields are set
// and that the storage backend is properly configured. Returns an error if the
// configuration is invalid.
func (c *RootConfig) Validate() error {
	c.Storage.Backend.Validate()

	if c.Services.Syslog.Protocol != "tcp" && c.Services.Syslog.Protocol != "udp" {
		return fmt.Errorf("invalid syslog protocol: %s", c.Services.Syslog.Protocol)
	}

	if c.Services.Syslog.Tls != nil && c.Services.Syslog.Protocol == "udp" {
		return fmt.Errorf("TLS is not supported for Syslog UDP protocol")
	}

	return nil
}

// Determines the storage backend type from the configuration and decodes the
// remaining body into the appropriate backend configuration struct. Returns an
// error if the backend is unsupported or if decoding fails.
func (c *StorageBackendConfig) Resolve() hcl.Diagnostics {
	switch c.Backend {
	case StorageBackendBadgerDb:
		config := DefaultStorageBackendBadgerDbConfig()

		diags := gohcl.DecodeBody(c.Body, nil, config)
		if diags.HasErrors() {
			c.BadgerDB = config
		}

		return diags

	case StorageBackendFoundationDb:
		panic("not implemented")

	default:
		return hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Unsupported storage backend",
				Detail:   fmt.Sprintf("The storage backend %q is not supported.", c.Backend),
			},
		}
	}
}

// Validates the storage configuration, ensuring that exactly one backend is
// configured. Returns an error if no backend is configured or if multiple
// backends are configured.
func (c *StorageBackendConfig) Validate() error {
	count := 0

	if c.BadgerDB != nil {
		count++
	}

	if count == 0 {
		return ErrNoStorageBackend
	} else if count > 1 {
		return ErrMultipleStorageBackends
	}

	return nil
}

// Load reads the configuration file at the given path, parses it, and
// decodes it into a RootConfig struct. It returns an error if parsing or
// decoding fails, or if the configuration is invalid.
func (cfg *RootConfig) Load(path string) error {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return diags
	}

	diags = diags.Extend(gohcl.DecodeBody(file.Body, nil, cfg))
	if diags.HasErrors() {
		return diags
	}

	if cfg.Storage == nil {
		return ErrNoStorageBackend
	}

	if cfg.Storage.Backend.Body != nil {
		diags = diags.Extend(cfg.Storage.Backend.Resolve())
		if diags.HasErrors() {
			return diags
		}
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	return nil
}

// Converts the parsed RootConfig into a server.Options struct, loading TLS
// certificates for the services that have TLS enabled.
func (cfg *RootConfig) AsServerOptions() (server.Options, error) {
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
	case StorageBackendBadgerDb:
		storageOptions = &server.BadgerDbStorageOptions{
			AuthDir:   cfg.Storage.Backend.BadgerDB.AuthDir,
			LogDir:    cfg.Storage.Backend.BadgerDB.LogDir,
			ConfigDir: cfg.Storage.Backend.BadgerDB.ConfigDir,
		}

	case StorageBackendFoundationDb:
		panic("not implemented")
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
