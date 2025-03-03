# Command Line Interface

```
Low-Code log management solution

Usage:
  flowg-server [flags]

Flags:
      --auth-dir string                   Path to the auth database directory (default "./data/auth")
      --config-dir string                 Path to the config directory (default "./data/config")
  -h, --help                              help for flowg-server
      --http-bind string                  Address to bind the HTTP server to (default ":5080")
      --http-tls                          Enable TLS for the HTTP server
      --http-tls-cert string              Path to the certificate file for the HTTPS server
      --http-tls-key string               Path to the certificate key file for the HTTPS server
      --log-dir string                    Path to the log database directory (default "./data/logs")
      --mgmt-bind string                  Address to bind the Management HTTP server to (default ":9113")
      --mgmt-tls                          Enable TLS for the Management HTTP server
      --mgmt-tls-cert string              Path to the certificate file for the Management HTTPS server
      --mgmt-tls-key string               Path to the certificate key file for the Management HTTPS server
      --syslog-allow-origin stringArray   Allowed origin (IP address or CIDR range) for Syslog server (default: all)
      --syslog-bind string                Address to bind the Syslog server to (default ":5514")
      --syslog-proto string               Protocol to use for the Syslog server (one of "tcp" or "udp") (default "udp")
      --syslog-tls                        Enable TLS for the Syslog server (requires protocol to be "tcp")
      --syslog-tls-auth                   Require clients to authenticate against the Syslog server with a client certificate
      --syslog-tls-cert string            Path to the certificate file for the Syslog server
      --syslog-tls-key string             Path to the certificate key file for the Syslog server
      --verbose                           Enable verbose logging
```
