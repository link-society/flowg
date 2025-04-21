---
sidebar_position: 2
---

# Automatic Cluster Formation Using Consul

One **FlowG** node can discover other **FlowG** nodes automatically if Consul is configured to be used.

In order to enable automatic cluster formation, have a instance of Consul running on the network.

Start the **FlowG** server with `--consul-url` flag set to the address of the Consul node. For example:

```bash
flowg-server \
  --syslog-proto="tcp" \
  --syslog-tls \
  --syslog-tls-cert="/path/to/cert.pem" \
  --syslog-tls-key="/path/to/cert.key" \
  --syslog-tls-auth \
  --consul-url="localhost:8500"
```

Alternatively, automatic cluster formation can also be enabled by using an environment variable. Set an environment variable `CONSUL_URL` to the address of the Consul node.
