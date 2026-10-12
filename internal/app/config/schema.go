// Package config defines the user-facing configuration of flowg-server: the
// schema of the HCL configuration file, its defaults (resolved from environment
// variables), and its validation.
package config

import (
	"errors"

	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
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

	BadgerDB     *StorageBackendBadgerDbConfig
	FoundationDB *StorageBackendFoundationDbConfig
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
	LogDir    string `hcl:"log_dir,optional"`
	ConfigDir string `hcl:"config_dir,optional"`
}

// Configuration for the FoundationDB storage backend, which is a distributed
// key-value store.
type StorageBackendFoundationDbConfig struct {
	ClusterFile      string `hcl:"cluster_file,optional"`
	ConnectionString string `hcl:"connection_string,optional"`
	KeySpace         string `hcl:"key_space,optional"`
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

// Returns the built-in default configuration for the server. It does not read
// the environment: environment variables are applied by [RootConfig.Load].
func DefaultConfig() *RootConfig {
	return &RootConfig{
		DemoMode: false,
		Logging: &LoggingConfig{
			Level:   "info",
			Verbose: false,
		},
		Services: &ServicesConfig{
			Http: &HttpConfig{
				BindAddress: ":5080",
				MountPath:   "/",
			},
			Management: &ManagementConfig{
				BindAddress: ":9113",
			},
			Syslog: &SyslogConfig{
				Bind:                  ":5514",
				Protocol:              "udp",
				InitialAllowedOrigins: []string{},
			},
		},
		Storage: &StorageConfig{
			Backend: &StorageBackendConfig{
				Backend:  StorageBackendBadgerDb,
				BadgerDB: DefaultStorageBackendBadgerDbConfig(),
			},
			Seed: &StorageSeedConfig{
				Auth: &StorageSeedAuthConfig{
					InitialUser:     "root",
					InitialPassword: "root",
				},
			},
		},
	}
}

// Returns the built-in default configuration for the BadgerDB storage backend.
func DefaultStorageBackendBadgerDbConfig() *StorageBackendBadgerDbConfig {
	return &StorageBackendBadgerDbConfig{
		AuthDir:   "./data/auth",
		LogDir:    "./data/logs",
		ConfigDir: "./data/config",
	}
}

// Returns the built-in default configuration for the FoundationDB storage
// backend.
func DefaultStorageBackendFoundationDbConfig() *StorageBackendFoundationDbConfig {
	return &StorageBackendFoundationDbConfig{
		ClusterFile: "/etc/foundationdb/fdb.cluster",
		KeySpace:    "flowg",
	}
}

// Validates the root configuration, ensuring that all required fields are set
// and that the storage backend is properly configured. Returns an error if the
// configuration is invalid.
func (c *RootConfig) Validate() error {
	if err := c.Storage.Backend.Validate(); err != nil {
		return err
	}

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
		if !diags.HasErrors() {
			c.BadgerDB = config
			c.FoundationDB = nil
		}

		return diags

	case StorageBackendFoundationDb:
		config := DefaultStorageBackendFoundationDbConfig()

		diags := gohcl.DecodeBody(c.Body, nil, config)
		if !diags.HasErrors() {
			c.FoundationDB = config
			c.BadgerDB = nil
		}

		return diags

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

// Validates the storage configuration, ensuring that the backend is supported
// and that exactly one backend is configured. Returns an error if the backend is
// unsupported, if no backend is configured or if multiple backends are
// configured.
func (c *StorageBackendConfig) Validate() error {
	if c.Backend != StorageBackendBadgerDb && c.Backend != StorageBackendFoundationDb {
		return fmt.Errorf("unsupported storage backend: %q", c.Backend)
	}

	count := 0

	if c.BadgerDB != nil {
		count++
	}

	if c.FoundationDB != nil {
		if c.FoundationDB.ClusterFile == "" && c.FoundationDB.ConnectionString == "" {
			return fmt.Errorf("FoundationDB backend requires either a cluster file or a connection string")
		}

		count++
	}

	if count == 0 {
		return ErrNoStorageBackend
	} else if count > 1 {
		return ErrMultipleStorageBackends
	}

	return nil
}

// Load completes the configuration with, by order of precedence, the
// environment variables that are set and the configuration file at the given
// path (skipped when the path is empty), the current values acting as
// defaults. It returns an error if parsing or decoding the file fails, or if
// the resulting configuration is invalid.
func (cfg *RootConfig) Load(path string) error {
	if path != "" {
		if err := cfg.loadFile(path); err != nil {
			return err
		}
	}

	cfg.applyEnv()

	return cfg.Validate()
}

// loadFile reads the configuration file at the given path, parses it, and
// decodes it over the current values. It returns an error if parsing or
// decoding fails.
func (cfg *RootConfig) loadFile(path string) error {
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

	return nil
}
