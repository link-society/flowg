# How To Build a Pipeline?

A pipeline is the entrypoint for logs in **FlowG**. Logs are ingested via the
REST API on a specific pipeline's endpoint.

As such, a pipeline flow will always have a root node named *Log Source*.

From this node, you are able to add the following type of nodes:

 - **Transform nodes:** Call a transformer to refine the log record and pass the
   result to the next nodes
 - **Switch nodes:** Pass the log record to the next nodes only if it matches
   the node's [filter](./filtering.md)
 - **Router nodes:** Store the log record into a stream

Using those nodes, a pipeline is able to parse, split, refine, enrich and route
log records to the database.
