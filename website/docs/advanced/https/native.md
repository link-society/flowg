---
sidebar_position: 3
---

# Configuring FlowG

## When deployed manually

Run **FlowG** with:

```bash
flowg serve \
  --auth-dir /var/lib/flowg/data/auth \
  --log-dir /var/lib/flowg/logs \
  --config-dir /var/lib/flowg/config \
  --http-bind 127.0.0.1:5080 \
  --http-tls \
  --http-tls-cert /etc/ssl/certs/logs.example.com.crt \
  --http-tls-key /etc/ssl/private/logs.example.com.key \
  --syslog-bind 127.0.0.1:5514
```

Or if using Certbot:

```bash
flowg serve \
  --auth-dir /var/lib/flowg/data/auth \
  --log-dir /var/lib/flowg/logs \
  --config-dir /var/lib/flowg/config \
  --http-bind 127.0.0.1:5080 \
  --http-tls \
  --http-tls-cert /etc/letsencrypt/live/logs.example.com/fullchain.pem \
  --http-tls-key /etc/letsencrypt/live/logs.example.com/privkey.pem \
  --syslog-bind 127.0.0.1:5514
```

## When deployed via Docker

Copy the certificates to a specific folder:

```bash
mkdir -p /opt/flowg/ssl
cp /etc/ssl/certs/logs.example.com.crt /opt/flowg/ssl/tls.crt
cp /etc/ssl/private/logs.example.com.key /opt/flowg/ssl/tls.key
```

Then run the Docker image with:

```
docker run \
  -p 5080:5080/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  -v /opt/flowg/ssl:/data/ssl \
  linksociety/flowg:latest serve \
    --http-tls \
    --http-tls-cert /data/ssl/tls.crt \
    --http-tls-key /data/ssl/tls.key
```

## When deployed on Kubernetes

Nothing to do, the Helm Chart configures **FlowG** automatically.
