---
sidebar_position: 3
---

# Configuring FlowG

## When deployed manually

Create a file named `/etc/flowg/config.hcl`:

```hcl
services {
  http {
    bind = "127.0.0.1:5080"

    tls {
      cert = "/etc/ssl/certs/logs.example.com.crt"
      key  = "/etc/ssl/private/logs.example.key"
    }
  }

  management {
    bind = "127.0.0.1:9113"

    tls {
      cert = "/etc/ssl/certs/mgmt.example.com.crt"
      key  = "/etc/ssl/private/mgmt.example.com.key"
    }
  }

  syslog {
    bind = "127.0.0.1:5514"
  }
}

storage {
  backend "badgerdb" {
    auth_dir = "/var/lib/flowg/auth"
    log_dir = "/var/lib/flowg/logs"
    config_dir = "/var/lib/flowg/config"
  }
}
```

Then, run **FlowG** with:

```bash
flowg-server --config /etc/flowg/config.hcl
```

Or if using Certbot:

```hcl

services {
  http {
    bind = "127.0.0.1:5080"

    tls {
      cert = "/etc/letsencrypt/live/logs.example.com/fullchain.pem"
      key  = "/etc/letsencrypt/live/logs.example.com/privkey.pem"
    }
  }

  management {
    bind = "127.0.0.1:9113"

    tls {
      cert = "/etc/letsencrypt/live/mgmt.example.com/fullchain.pem"
      key  = "/etc/letsencrypt/live/mgmt.example.com/privkey.key"
    }
  }

  syslog {
    bind = "127.0.0.1:5514"
  }
}

storage {
  backend "badgerdb" {
    auth_dir = "/var/lib/flowg/auth"
    log_dir = "/var/lib/flowg/logs"
    config_dir = "/var/lib/flowg/config"
  }
}
```

## When deployed via Docker

Copy the certificates to a specific folder:

```bash
mkdir -p /opt/flowg/ssl
cp /etc/ssl/certs/logs.example.com.crt /opt/flowg/ssl/tls.crt
cp /etc/ssl/private/logs.example.com.key /opt/flowg/ssl/tls.key

cp /etc/ssl/certs/mgmt.example.com.crt /opt/flowg/ssl/tls-mgmt.crt
cp /etc/ssl/private/mgmt.example.com.key /opt/flowg/ssl/tls-mgmt.key
```

Then run the Docker image with:

```
docker run \
  -p 5080:5080/tcp \
  -p 9113:9113/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  -v /opt/flowg/ssl:/data/ssl \
  -e FLOWG_HTTP_TLS_ENABLED=true \
  -e FLOWG_HTTP_TLS_CERT=/data/ssl/tls.crt \
  -e FLOWG_HTTP_TLS_KEY=/data/ssl/tls.key \
  -e FLOWG_MGMT_TLS_ENABLED=true \
  -e FLOWG_MGMT_TLS_CERT=/data/ssl/tls-mgmt.crt \
  -e FLOWG_MGMT_TLS_KEY=/data/ssl/tls-mgmt.key \
  linksociety/flowg:latest
```

## When deployed on Kubernetes

Nothing to do, the Helm Chart configures **FlowG** automatically.
