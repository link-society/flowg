---
sidebar_position: 6
---

# Configuration File

The configuration file format is the
[Hashicorp Configuration Language](https://github.com/hashicorp/hcl).

## Example

Here is an exhaustive example of the configuration file:

```hcl
# env: FLOWG_DEMO_MODE (default: false)
demo = false

logging {
  # env: FLOWG_VERBOSE (default: false)
  verbose = false

  # env: FLOWG_LOGLEVEL (default: "info")
  level = "info"
}

services {
  http {
    # env: FLOWG_HTTP_BIND_ADDRESS (default: ":5080")
    bind = ":5080"

    # env: FLOWG_HTTP_MOUNT_PATH (default: "/")
    mount = "/"

    # env: FLOWG_HTTP_TLS_ENABLED (default: false)
    tls {
      # env: FLOWG_HTTP_TLS_CERT (default: "")
      cert = "/run/secrets/tls.crt"

      # env: FLOWG_HTTP_TLS_KEY (default: "")
      key  = "/run/secrets/tls.key"
    }
  }

  management {
    # env: FLOWG_MGMT_BIND_ADDRESS (default: ":9113")
    bind = ":9113"

    # env: FLOWG_MGMT_TLS_ENABLED (default: false)
    tls {
      # env: FLOWG_MGMT_TLS_CERT (default: "")
      cert = "/run/secrets/tls.management.crt"

      # env: FLOWG_MGMT_TLS_KEY (default: "")
      key  = "/run/secrets/tls.management.key"
    }
  }

  syslog {
    # env: FLOWG_SYSLOG_BIND_ADDRESS (default: ":5514")
    bind = ":5514"

    # env: FLOWG_SYSLOG_PROTOCOL (default: "udp")
    protocol = "tcp"

    # env: FLOWG_SYSLOG_INITIAL_ALLOWED_ORIGINS (default: [])
    initial_allowed_origins = []

    # env: FLOWG_SYSLOG_TLS_ENABLED (default: false)
    tls {
      # env: FLOWG_SYSLOG_TLS_CERT (default: "")
      cert = "/run/secrets/tls.syslog.crt"

      # env: FLOWG_SYSLOG_TLS_KEY (default: "")
      key  = "/run/secrets/tls.syslog.key"

      # env: FLOWG_SYSLOG_TLS_AUTH (default: false)
      auth = true
    }
  }
}

storage {
  # env: FLOWG_STORAGE_BACKEND (default: "badgerdb")
  backend "badgerdb" {
    # env: FLOWG_BADGER_AUTH_DIR (default: "./data/auth")
    auth_dir = "./data/auth"

    # env: FLOWG_BADGER_LOG_DIR (default: "./data/logs")
    log_dir = "./data/logs"

    # env: FLOWG_BADGER_CONFIG_DIR (default: "./data/config")
    config_dir = "./data/config"
  }

  # env: FLOWG_STORAGE_BACKEND
  backend "foundationdb" {
    # env: FLOWG_FOUNDATIONDB_CLUSTER_FILE (default: "/etc/foundationdb/fdb.cluster")
    cluster_file = "/etc/foundationdb/fdb.cluster"

    # env: FLOWG_FOUNDATIONDB_CONNECTION_STRING (default: "")
    # if set, has priority over `cluster_file`
    connection_string = ""

    # env: FLOWG_FOUNDATIONDB_KEY_SPACE (default: "flowg")
    key_space = ""
  }

  # NB: only one storage backend section is allowed

  seed {
    auth {
      # env: FLOWG_AUTH_INITIAL_USER (default: "root")
      initial_user = "root"

      # env: FLOWG_AUTH_INITIAL_PASSWORD (default: "root")
      initial_password = "root"

      # env: FLOWG_AUTH_RESET_USER (default: "")
      reset_user = ""

      # env: FLOWG_AUTH_RESET_PASSWORD (default: "")
      reset_password = ""
    }
  }
}

security {
  # env: FLOWG_SECRET_KEY (default: "", a random key is generated)
  secret_key = "a-long-random-secret"
}
```

## Priority

Settings specified in environment variables have the highest priority. If
unset (or empty), their value is read from the configuration file. If missing,
a default value is provided.

Setting `FLOWG_STORAGE_BACKEND` to a backend other than the one configured in
the file discards the settings of the file's backend: the selected backend
starts from its default values, overridden by its own environment variables.

Setting a `FLOWG_*_TLS_ENABLED` variable to `false` disables TLS for the
service, even if the configuration file has a `tls` block.

## Security

JSON Web Tokens, one of the mechanisms used for authentication against the API,
are signed using the key provided by `security.secret_key` (or the environment
variable `FLOWG_SECRET_KEY`).

If unset, a random key is generated at startup: issued tokens are then
invalidated whenever the server restarts. In a cluster, every node must share
the same key.
