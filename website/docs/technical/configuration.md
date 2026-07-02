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
```

## Priority

Settings specified in the configuration file have the highest priority. If
missing, their value is read from environment variables. If unset, a default
value is provided.

## Security

JSON Web Tokens, one of the mechanisms used for authentication against the API,
are signed using the key provided by the environment variable
`FLOWG_SECRET_KEY`.

If the environment variable is unset, a random key will be used.
