---
sidebar_position: 3
---

# Configuring FlowG

## When deployed manually

Run **FlowG** with:

```bash
flowg-server \
  --auth-dir /var/lib/flowg/data/auth \
  --log-dir /var/lib/flowg/logs \
  --config-dir /var/lib/flowg/config \
  --http-bind 127.0.0.1:5080 \
  --http-tls \
  --http-tls-cert /etc/ssl/certs/logs.example.com.crt \
  --http-tls-key /etc/ssl/private/logs.example.com.key \
  --mgmt-bind 127.0.0.1:9113 \
  --mgmt-tls \
  --mgmt-tls-cert /etc/ssl/certs/mgmt.example.com.crt \
  --mgmt-tls-key /etc/ssl/private/mgmt.example.com.key \
  --syslog-bind 127.0.0.1:5514
```

Or if using Certbot:

```bash
flowg-server \
  --auth-dir /var/lib/flowg/data/auth \
  --log-dir /var/lib/flowg/logs \
  --config-dir /var/lib/flowg/config \
  --http-bind 127.0.0.1:5080 \
  --http-tls \
  --http-tls-cert /etc/letsencrypt/live/logs.example.com/fullchain.pem \
  --http-tls-key /etc/letsencrypt/live/logs.example.com/privkey.pem \
  --mgmt-bind 127.0.0.1:9113 \
  --mgmt-tls \
  --mgmt-tls-cert /etc/letsencrypt/live/mgmt.example.com/fullchain.pem \
  --mgmt-tls-key /etc/letsencrypt/live/mgmt.example.com/privkey.key \
  --syslog-bind 127.0.0.1:5514
```

## When deployed via Docker

Copy the certificates to a specific folder:

```bash
mkdir -p /opt/flowg/ssl
cp /etc/ssl/certs/logs.example.com.crt /opt/flowg/ssl/tls.crt
cp /etc/ssl/private/logs.example.com.key /opt/flowg/ssl/tls.key

cp /etc/ssl/certs/mgmt.example.com.crt /opt/flowg/tls-mgmt.crt
cp /etc/ssl/private/mgmt.example.com.key /opt/flowg/tls-mgmt.key
```

Then run the Docker image with:

```
docker run \
  -p 5080:5080/tcp \
  -p 9113:9113/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  -v /opt/flowg/ssl:/data/ssl \
  linksociety/flowg:latest serve \
    --http-tls \
    --http-tls-cert /data/ssl/tls.crt \
    --http-tls-key /data/ssl/tls.key \
    --mgmt-tls \
    --mgmt-tls-cert /data/ssl/tls-mgmt.crt \
    --mgmt-tls-key /data/ssl/tls-mgmt.key
```

## When deployed on Kubernetes

Nothing to do, the Helm Chart configures **FlowG** automatically.
