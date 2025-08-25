---
title: ElasticSearch API compatibility
description: FlowG v0.45.0 introduces API compatibility with ElasticSearch
slug: elasticsearch-compatibility
authors: linkdd
tags: [release, elastcisearch, compat]
---

:tada:
**[FlowG v0.45.0](https://github.com/link-society/flowg/releases/tag/v0.45.0)**
has been released with partial support for the ElasticSearch API!

<!-- truncate -->

## Introduction

**FlowG**'s goals have always been: interoperability and ease of use. Many
applications already use the "ELK" stack:

 - *ElasticSearch* for indexing
 - *Logstash* for aggregation
 - *Kibana* for viewing

Usually, logs are sent to *Logstash* via Syslog, which then forwards them to
ElasticSearch for storage and indexing.

**FlowG** already could be set up as a drop-in replacement for *Logstash* thanks
to its Syslog Server endpoint, and its multitude of forwarders.

But with the latest release, we're taking things up a notch with the partial
support for ElasticSearch API.

## What does this mean?

**FlowG** exposes on the `/api/v1/middlewares/elastic` endpoint an API
compatible with ElasticSearch, allowing you to plug your existing application,
which uses the ElasticSearch client libraries, into **FlowG**, without changing
your code.

> :warning: **NB:** The support is only partial.

At the moment, only 2 operations are supported:

 - `HEAD /{index}`: check if an index exists
 - `POST /{index}/_doc`: index a document

And only HTTP Basic authentication is supported.

Indexes map to FlowG pipelines, the following request would send the document
through the `default` pipeline:

```
POST /api/v1/middlewares/elastic/default/_doc
Authorization: Basic ...
{"foo": {"bar": "baz}}
```

> **NB:** In FlowG, the datamodel is flat, so the document is flattened by the
> API before handing it to the pipeline, this would be equivalent to:

```json
{"foo.bar": "baz"}
```

## What now? The roadmap

More operations might be added later on to the "compatibility API" to smooth out
the integration of **FlowG** into your existing infrastructure, though the goal
is not to have 100% feature parity (only the subset that makes sense to support
in **FlowG**).

But also, more APIs might come in later, depending on user feedback/requests.

## Summary

As of **FlowG v0.45.0**, logs can be ingested via:

 - HTTP as text data (one log per line) on `/api/v1/pipelines/{PIPELINE}/logs/text`
 - HTTP as JSON data on `/api/v1/pipelines/{PIPELINE}/logs/struct`
 - HTTP as OpenTelemetry data on `/api/v1/pipelines/{PIPELINE}/logs/otlp`
 - Syslog protocol
 - ElasticSearch REST API on `/api/v1/middlewares/elastic/{PIPELINE}/_doc`

[See the release notes](https://github.com/link-society/flowg/releases/tag/v0.45.0)