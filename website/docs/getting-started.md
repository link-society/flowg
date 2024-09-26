---
sidebar_position: 1
---

# Getting Started

**FlowG** is a log processing platform that lets you ingest, transform, and
query logs using a visual pipeline builder. It handles structured logs without
requiring predefined schemas.

**FlowG**'s primary goal is to make log refinement and routing as easy as
possible. As such, it relies on [React Flow](https://reactflow.dev) to help you
build such pipelines with as little code as possible. For the actual code part,
we use the [Vector Remap Language](https://vector.dev/docs/reference/vrl/),
which gives a solid base for log refinement.

It aims to replace tools like [Logstash](https://www.elastic.co/logstash) by
integrating the feature right in the solution.

It also leverages [BadgerDB](https://dgraph.io/docs/badger/) which is a battle
tested Key/Value database, with the right feature set to easily support indexing
dynamically structured logs, as well as log compression.

## Installation

### Using Docker


```bash
docker run \
  -p 5080:5080/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  linksociety/flowg:latest serve
```

### Using Kubernetes

```bash
git clone https://github.com/link-society/flowg
cd flowg
helm install flowg ./k8s/charts/flowg \
  --create-namespace \
  --namespace flowg-system \
  --wait
kubectl port-forward svc/flowg 5080:5080 --namespace flowg-system
```

This will automatically deploy [Fluentd](https://www.fluentd.org) alongside
**FlowG** in order to collect logs from all pods.

## Next Steps

Connect to [http://localhost:5080](http://localhost:5080) with the default
credentials `root` / `root`.
