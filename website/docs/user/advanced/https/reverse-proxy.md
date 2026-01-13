---
sidebar_position: 2
---

# Configuring a Reverse Proxy

By default, **FlowG** listens on port `5080`. To run on port `80` or `443`, it
would require root privileges, which is not necessarily a good idea.

Using a Reverse Proxy solves that issue, **FlowG** can still run on port `5080`.

> **NB:** By default, the Management interface listens on port `9113`. It is not
> recommended to expose it publicly.

## Using Apache2

In `/etc/apache2/sites-enabled/flowg.conf`:

```apacheconf
<VirtualHost *:80>
  ServerName logs.example.com

  Redirect permanent / https://logs.example.com/

  ErrorLog ${APACHE_LOG_DIR}/flowg/error.log
  CustomLog ${APACHE_LOG_DIR}/flowg/access.log combined
</VirtualHost>

<VirtualHost *:443>
  ServerName logs.example.com

  SSLEngine On
  SSLCertificateFile /etc/ssl/certs/logs.example.com.crt
  SSLCertificateKeyFile /etc/ssl/private/logs.example.com.key

  ## Or if using Certbot
  # SSLCertificateFile /etc/letsencrypt/live/logs.example.com/fullchain.pem
  # SSLCertificateKeyFile /etc/letsencrypt/live/logs.example.com/privkey.pem

  ProxyPass / http://127.0.0.1:5080/
  ProxyPassReverse / http://127.0.0.1:5080/
  ProxyPreserveHost On

  ErrorLog ${APACHE_LOG_DIR}/flowg/error-tls.log
  CustomLog ${APACHE_LOG_DIR}/flowg/access-tls.log combined
</VirtualHost>
```

## Using NGINX

In `/etc/nginx/sites-enabled/flowg.conf`:

```nginx
server {
  server_name logs.example.com;

  listen 80;
  return 301 https://$host$request_uri;
}

server {
  server_name logs.example.com;

  listen 443 ssl;

  ssl_certificate /etc/ssl/certs/logs.example.com.crt;
  ssl_certificate_key /etc/ssl/private/logs.example.com.key;

  ## Or if using Certbot
  # ssl_certificate /etc/letsencrypt/live/logs.example.com/fullchain.pem;
  # ssl_certificate_key /etc/letsencrypt/live/logs.example.com/privkey.pem;

  location / {
    proxy_pass http://127.0.0.1:5080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
  }

  access_log /var/log/nginx/flowg/access.log;
  error_log /var/log/nginx/flowg/error.log;
}
```

## Using Caddy

In `/etc/caddy/Caddyfile`:

```hcl
logs.example.com {
  reverse_proxy http://127.0.0.1:5080

  log {
    output_file /var/log/caddy/flowg/access.log
    level info
  }
}
```
