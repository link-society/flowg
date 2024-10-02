---
sidebar_position: 1
---

# Creating the server certificate

The first step in setting up HTTPS is to obtain the server certificate. This
page will cover 3 ways to create said certificate, choose the one that best fit
your use case.

In this guide, we will assume you are exposing **FlowG** via the
`logs.example.com` hostname.

> **NB:** If using Caddy as a reverse proxy, this step can be skipped as it
> supports automatic ACME challenges.

## Using a self-signed certificate

A self-signed certificate is often used for pre-production/staging environments,
or development environments. They are not really suitable for production, but
it's still a good way to start.

First, generate a private key:

```bash
openssl genpkey -algorithm RSA -out server.key -aes256
```

> **NB:** If you don't want to use a passphrase, omit the `-aes256` flag

Then, in a configuration file named `openssl.cnf`:

```ini
[ req ]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[ req_distinguished_name ]
C = US
ST = State
L = City
O = Organization
OU = Organizational Unit
CN = logs.example.com
emailAddress = admin@logs.example.com

[ v3_req ]
```

> **NB:** Modify the fields to your convenience.

Finally, generate the certificate:

```bash
openssl req -x509 -nodes \
  -days 365 \
  -key server.key \
  -out server.crt \
  -config openssl.cnf
```

And install the certificates:

```bash
sudo mv server.key /etc/ssl/private/logs.example.com.key
sudo chmod 600 /etc/ssl/private/logs.example.com.key
sudo chown root:root /etc/ssl/private/logs.example.com.key

sudo mv server.crt /etc/ssl/certs/logs.example.com.crt
sudo chmod 644 /etc/ssl/certs/logs.example.com.crt
sudo chown root:root /etc/ssl/certs/logs.example.com.crt
```

## Using certbot

First, install `certbot`:

```bash
apt install certbot
```

Then run the following command:

```bash
sudo certbot certonly --standalone -d logs.example.com
```

If you were running a webserver, for example NGINX, you need to stop it during
the challenge:

```bash
sudo certbot certonly --standalone \
  --pre-hook "systemctl stop nginx" \
  --post-hook "systemctl start nginx" \
  -d logs.example.com
```

To renew the certificate:

```bash
sudo certbot renew --standalone --cert-name logs.example.com
```

Or with hooks:

```bash
sudo certbot renew --standalone \
  --pre-hook "systemctl stop nginx" \
  --post-hook "systemctl start nginx" \
  --cert-name logs.example.com
```

You can automate the renewal by adding the renew command to your *crontab*:

```bash
cat <<EOF | sudo tee /etc/cron.daily/certbot-flowg
#!/bin/bash
certbot renew --standalone --quiet --cert-name logs.example.com

## Or with hooks:
# certbot renew --standalone --quiet \
#   --pre-hook "systemctl stop nginx" \
#   --post-hook "systemctl start nginx" \
#   --cert-name logs.example.com
EOF

sudo chmox +x /etc/cron.daily/certbot-flowg
```

Certificates will be stored in `/etc/letsencrypt/live/logs.example.com/`.

## Using Cert-Manager (on Kubernetes)

If you deployed **FlowG** on Kubernetes, using the Helm Chart, you can rely on
[Cert-Manager](https://cert-manager.io) to issue and renew the certificate.

You will need to redeploy the Helm Chart with the following values:

```yaml
---

# ...

flowg:
  # ...

  https:
    enabled: true

    certificateFrom:
      ## Manually issued certificate stored in a `Secret`:
      # secretRef:
      #   ## Expected keys: `tls.crt` & `tls.key`
      #   name: logs-example-com-tls

      ## Automatically issued with Cert-Manager
      certmanager:
        commonName: logs.example.com
        issuerRef:
          name: ca-issuer      # adjust to your setup
          kind: ClusterIssuer  # adjust to your setup

# ...
```

> **NB:** The Helm Chart will not deploy an `Issuer` or `ClusterIssuer`
> resource, it must exist prior to the Helm Chart's deployment.
