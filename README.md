# FlowG

FlowG is a log management platform that lets you ingest, transform, and query
logs using a visual pipeline builder. It handles structured logs without
requiring predefined schemas and relies on
[BadgerDB](https://dgraph.io/docs/badger/) as its storage backend.

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

## :warning: Disclaimer

This project is in a very early stage of development, as such:

 - it has not been tested in production
 - no benchmark available
 - not even an automated test suite

Those points will be improved on given enough time.

## :construction: Build

**Requirements:**

 - Go 1.22+
 - C Toolchain
 - Rust and Cargo (edition 2021)
 - NodeJS

```
$ go install github.com/go-task/task/v3/cmd/task@latest
$ task build
```

Then run:

```
$ ./bin/flowg \
  -db ./data/logs \
  -config ./data/config  \
  -bind 127.0.0.1:5080
```

Then you can access:

 - the WebUI at http://localhost:5080
 - the API documentation at http://localhost:5080/api/docs

## :whale: Build the Docker Image

```
$ task docker:build
```

This will build `linksociety/flowg:latest` locally.

Then run with:

```
$ docker run -p 5080:5080 -v ./data:/data linksociety/flowg:latest
```

## :memo: License

This software is released under the terms of the [MIT License](./LICENSE.txt)
