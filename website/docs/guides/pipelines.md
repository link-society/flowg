---
sidebar_position: 5
---

# How To Build a Pipeline?

A pipeline is the entrypoint for logs in **FlowG**. Logs can be ingested via:

 - the REST API on a specific pipeline's endpoint
 - the Syslog Server endpoint (UDP)

As such, a pipeline flow will always have 2 root nodes:

 - `DIRECT`: for logs ingested via the pipeline's API endpoint
 - `SYSLOG`: for logs received via the Syslog endpoint

From those nodes, you are able to add the following type of nodes:

 - **Transform nodes:** Call a transformer to refine the log record and pass the
   result to the next nodes
 - **Switch nodes:** Pass the log record to the next nodes only if it matches
   the node's [filter](/docs/guides/filtering)
 - **Pipeline nodes:** Pass the log record to another pipeline
 - **Alert nodes:** Send the log to an [Alert webhook](/docs/guides/alerting)
 - **Router nodes:** Store the log record into a stream

Using those nodes, a pipeline is able to parse, split, refine, enrich and route
log records to the database.
