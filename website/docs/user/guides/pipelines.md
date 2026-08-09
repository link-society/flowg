---
sidebar_position: 4
---

# How To Build a Pipeline?

## Introduction

A pipeline is the entrypoint for logs in **FlowG**. Logs can be ingested via:

 - the REST API on a specific pipeline's endpoint
 - the Syslog Server endpoint (UDP, TCP, or TCP+TLS)

As such, a pipeline flow will always have 2 root nodes:

 - `DIRECT`: for logs ingested via the pipeline's API endpoints
 - `SYSLOG`: for logs received via the Syslog endpoint

From those nodes, you are able to add the following types of node:

 - **Transform nodes:** Call a transformer to refine the log record and pass the
   result to the next nodes
 - **Switch nodes:** Pass the log record to the next nodes only if it matches
   the node's [filter](/docs/user/guides/filtering)
 - **Metric nodes:** Count the number of log records that goes through them
 - **Pipeline nodes:** Pass the log record to another pipeline
 - **Forward nodes:** Send the log to a third-party service
 - **Router nodes:** Store the log record into a stream

Using those nodes, a pipeline is able to parse, split, refine, enrich and route
log records to the database.

For more information, consult the [Technical Documentation](/docs/technical/pipelines).

## Prometheus Exporter

Every pipeline provides a Prometheus Exporter that can be scrapped:

```
GET /api/v1/pipelines/MYPIPELINE/metrics
```

```
# HELP node_example1 Number of logs measured since startup
# TYPE node_example1 counter
node_example1 23
# HELP node_example2 Number of logs measured since startup
# TYPE node_example2 counter
node_example2 42
# HELP pipeline_logs_total Total number of logs processed by the pipeline since startup
# TYPE pipeline_logs_total counter
pipeline_logs_total 65
```
