// Package config defines the user-facing configuration of flowg-server: the
// schema of the HCL configuration file, its defaults (declared by the `default`
// struct tags), its environment overlay, and its validation.
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
	DemoMode bool `hcl:"demo,optional" default:"false"`

	Logging  *LoggingConfig  `hcl:"logging,block" default:"{}"`
	Services *ServicesConfig `hcl:"services,block" default:"{}"`
	Storage  *StorageConfig  `hcl:"storage,block" default:"{}"`
}

// Define logging settings
type LoggingConfig struct {
	Level   string `hcl:"level,optional" default:"info"`
	Verbose bool   `hcl:"verbose,optional" default:"false"`
}

// Group together the configuration for the services provided by the server.
// Each service needs to bind to a specific address and port, and may have TLS
// support.
type ServicesConfig struct {
	Http       *HttpConfig       `hcl:"http,block" default:"{}"`
	Management *ManagementConfig `hcl:"management,block" default:"{}"`
	Syslog     *SyslogConfig     `hcl:"syslog,block" default:"{}"`
}

// Configuration for the "HTTP" service, which provides the REST API and web UI.
type HttpConfig struct {
	BindAddress string         `hcl:"bind,optional" default:":5080"`
	MountPath   string         `hcl:"mount,optional" default:"/"`
	Tls         *HttpTlsConfig `hcl:"tls,block"`
}

// Configuration for the "Management" service, which provides the health check
// and metrics endpoints.
type ManagementConfig struct {
	BindAddress string         `hcl:"bind,optional" default:":9113"`
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
	Bind                  string           `hcl:"bind,optional" default:":5514"`
	Protocol              string           `hcl:"protocol,optional" default:"udp"`
	InitialAllowedOrigins []string         `hcl:"initial_allowed_origins,optional" default:""`
	Tls                   *SyslogTlsConfig `hcl:"tls,block"`
}

// Configuration for TLS settings for the Syslog service.
type SyslogTlsConfig struct {
	Cert        string `hcl:"cert"`
	Key         string `hcl:"key"`
	AuthEnabled bool   `hcl:"auth,optional" default:"false"`
}

// StorageConfig groups together settings related to storage.
type StorageConfig struct {
	Backend *StorageBackendConfig `hcl:"backend,block" default:"{}"`
	Seed    *StorageSeedConfig    `hcl:"seed,block" default:"{}"`
}

// Configuration for the storage backend used by the server. It supports
// multiple backends, but only one can be configured at a time: the settings of
// the selected backend are held by the matching field, the others being nil.
type StorageBackendConfig struct {
	Backend string   `hcl:"backend,label" default:"badgerdb"`
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
	AuthDir   string `hcl:"auth_dir,optional" default:"./data/auth"`
	LogDir    string `hcl:"log_dir,optional" default:"./data/logs"`
	ConfigDir string `hcl:"config_dir,optional" default:"./data/config"`
}

// Configuration for the FoundationDB storage backend, which is a distributed
// key-value store.
type StorageBackendFoundationDbConfig struct {
	ClusterFile      string `hcl:"cluster_file,optional" default:"/etc/foundationdb/fdb.cluster"`
	ConnectionString string `hcl:"connection_string,optional" default:""`
	KeySpace         string `hcl:"key_space,optional" default:"flowg"`
}

// Configuration for seeding the storage backend with initial data. This is
// useful for setting up default users, roles, and other necessary data during
// the first run of the server.
type StorageSeedConfig struct {
	Auth *StorageSeedAuthConfig `hcl:"auth,block" default:"{}"`
}

// Configuration for seeding the authentication storage with an initial user
// and password. This is useful for bootstrapping the server with a default
// admin account.
type StorageSeedAuthConfig struct {
	InitialUser     string `hcl:"initial_user" default:"root"`
	InitialPassword string `hcl:"initial_password" default:"root"`

	ResetUser     string `hcl:"reset_user,optional" default:""`
	ResetPassword string `hcl:"reset_password,optional" default:""`
}

// Returns the built-in default configuration for the server, as declared by
// the `default` struct tags of the schema. It does not read the environment:
// environment variables are applied by [RootConfig.Load].
func DefaultConfig() *RootConfig {
	cfg := newWithDefaults[RootConfig]()
	cfg.Storage.Backend.selectBackend(cfg.Storage.Backend.Backend)
	return cfg
}

// Returns the built-in default configuration for the BadgerDB storage backend.
func DefaultStorageBackendBadgerDbConfig() *StorageBackendBadgerDbConfig {
	return newWithDefaults[StorageBackendBadgerDbConfig]()
}

// Returns the built-in default configuration for the FoundationDB storage
// backend.
func DefaultStorageBackendFoundationDbConfig() *StorageBackendFoundationDbConfig {
	return newWithDefaults[StorageBackendFoundationDbConfig]()
}

// selectBackend selects the named backend, with its default settings, and
// drops the settings of any other backend. An unsupported name leaves no
// backend settings, which [StorageBackendConfig.Validate] reports.
func (c *StorageBackendConfig) selectBackend(name string) {
	c.Backend = name
	c.BadgerDB = nil
	c.FoundationDB = nil

	switch name {
	case StorageBackendBadgerDb:
		c.BadgerDB = DefaultStorageBackendBadgerDbConfig()

	case StorageBackendFoundationDb:
		c.FoundationDB = DefaultStorageBackendFoundationDbConfig()
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
