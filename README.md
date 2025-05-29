<div align="center">

  ![Logo](./website/static/img/logo.png)

  # FlowG

  [![Build Status](https://github.com/link-society/flowg/actions/workflows/ci-package.yml/badge.svg)](https://github.com/link-society/flowg/actions/workflows/ci-package.yml)
  [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=link-society_flowg&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=link-society_flowg)
  [![Discord](https://img.shields.io/discord/1377785473363345540)](https://discord.gg/r6fvz7K3)

</div>

FlowG is a log management platform that lets you ingest, transform, and query
logs using a visual pipeline builder. It handles structured logs without
requiring predefined schemas and relies on
[BadgerDB](https://dgraph.io/docs/badger/) as its storage backend.

For more informations, consult the
[documentation](https://link-society.github.io/flowg/docs/).

You wish to contribute? Please read our [guidelines](./CONTRIBUTING.md) first.

## :sparkles: Key Features

 - **Pipelines:** Build and manage log processing pipelines visually, with
   support for transormations with [VRL](https://vector.dev/docs/reference/vrl/)
   scripts and conditional routing to dedicated streams
 - **Dynamic Logs:** Manage logs with varying structures without needing
   predefined schemas
 - **Interoperability:** Configuration and log ingestion can entirely be done
   via the REST API
 - **Observability:** Query and visualize logs in real-time in the integrated
   Web interface

![pipeline screenshot](./docs/screenshots/pipelines.png)

![stream screenshot](./docs/screenshots/streams.png)

## :question: Why should I use this?

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

## :zap: Benchmark

Run with (see *Build* section bellow):

```bash
export BENCHMARK_ITERATIONS="1000000"
export BENCHMARK_JOBS="24"  # usually, number of CPUs
task docker:build test:bench
```

System:

 - **CPU:** 13th Gen Intel(R) Core(TM) i7-13700K
 - **RAM:** 64GB

Result:

```
Requests sent: 1000000
Total time:    266.75s
Rate:          3748.76 req/s
```

## :building_construction: Build

**Requirements:**

 - Go 1.22+
 - C Toolchain
 - Rust and Cargo (edition 2021)
 - NodeJS

```bash
go install github.com/go-task/task/v3/cmd/task@latest
task build:all
```

Then, start the server with:

```bash
./bin/flowg-server \
  --auth-dir ./data/auth \
  --log-dir ./data/logs \
  --config-dir ./data/config \
  --http-bind 127.0.0.1:5080 \
  --mgmt-bind 127.0.0.1:9113 \
  --syslog-bind 127.0.0.1:5514
```

Now, you can access:

 - the WebUI at http://localhost:5080
 - the API documentation at http://localhost:5080/api/docs
 - the Syslog Server at [udp://localhost:5514]()
 - the Management API at [http://localhost:9113](http://localhost:9113):
    - `/health`: health check, always return `200 OK`
    - `/metrics`: Prometheus Exporter

A default user `root` (password: `root`) and a default pipeline are bootsrapped
if no configuration exists during startup.

To build the Docker image `linksociety/flowg:latest` locally:

```bash
task docker:build
```

## :rocket: Deploy

Using Docker:

```bash
docker run \
  -p 5080:5080/tcp \
  -p 9113:9113/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  linksociety/flowg:latest
```

Using Kubernetes (and Helm):

```bash
helm install flowg ./k8s/charts/flowg -n flowg-system --create-namespace
```

## :memo: License

This software is released under the terms of the [MIT License](./LICENSE.txt)
