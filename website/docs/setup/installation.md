---
sidebar_position: 1
---

# Installation

### From Sources

**Requirements:**

 - Go 1.22+
 - C Toolchain
 - Rust and Cargo (edition 2021)
 - NodeJS

First fetch the source code:

```bash
git clone https://github.com/link-society/flowg
cd flowg
```

Then run the build system:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
task build
sudo install -m 755 ./bin/flowg /usr/local/bin/flowg
```

Then, start the server with:

```bash
flowg serve \
  --auth-dir /var/lib/flowg/data/auth \
  --log-dir /var/lib/flowg/logs \
  --config-dir /var/lib/flowg/config \
  --http-bind 127.0.0.1:5080 \
  --syslog-bind 127.0.0.1:5514
```

> **NB:** All the options are optional and default to the values shown above.

For more informations about the command line interface, please consult
[this document](https://github.com/link-society/flowg/blob/main/docs/cli.md).

### Using Docker

```bash
docker run \
  -p 5080:5080/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  linksociety/flowg:latest serve
```

### Using Kubernetes

First fetch the source code:

```bash
git clone https://github.com/link-society/flowg
cd flowg
```

Then deploy the Helm chart:

```bash
helm install flowg ./k8s/charts/flowg \
  --create-namespace \
  --namespace flowg-system \
  --wait
```

This will automatically deploy [Fluentd](https://www.fluentd.org) alongside
**FlowG** in order to collect logs from all pods.

Once deployed, you can forward the port of the WebUI and API to your localhost:

```bash
kubectl port-forward svc/flowg 5080:5080 --namespace flowg-system
```

## Next Steps

Once deployed, **FlowG** creates a default pipeline and a default account with
the credentials `root` / `root`.

You now have access to:

 - The WebUI at [http://localhost:5080](http://localhost:5080)
 - The API documentation at [http://localhost:5080/api/docs](https://localhost:5080/api/docs)
 - The Syslog endpoint at [udp://localhost:5514](udp://localhost:5514)