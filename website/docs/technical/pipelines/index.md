---
sidebar_position: 1
---

# Pipelines

A pipeline is a flow-based graph through which log records are processed.

Logs can be ingested via the HTTP API, or via the embedded Syslog Server. As
such, every pipeline has 2 entrypoints, but can have many termination nodes.

## Data Model

```mermaid
erDiagram
    Flow {
        int version "For automatic migration"
        Node[] nodes
        Edge[] edges
    }
    Node {
        str id UK "Either builtin, or with an UUID"
        vec2 position "Position in the canvas"
        str type "Node type, determine the schema of `data`"
        any data "Node configuration"
    }
    Edge {
        str id UK "Composed of the source/target node IDs"
        str sourceHandle "Unused"
        str source "ID of the source node"
        str target "ID of the target node"
    }
   Flow ||--o{ Node : contains
   Flow ||--o{ Edge : contains
   Edge ||--|| Node : references
```
