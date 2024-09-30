---
sidebar_position: 1
---

# Introduction

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
