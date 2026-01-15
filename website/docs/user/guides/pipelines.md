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
 - **Pipeline nodes:** Pass the log record to another pipeline
 - **Forward nodes:** Send the log to a third-party service
 - **Router nodes:** Store the log record into a stream

Using those nodes, a pipeline is able to parse, split, refine, enrich and route
log records to the database.

For more information, consult the [Technical Documentation](/docs/technical/pipelines).
