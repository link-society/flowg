---
sidebar_position: 3
---

# Architecture

## Component Overview

A node is built from a few cooperating components layered on top of FlowG's HTTP
management interface:

```mermaid
flowchart TB
  subgraph Node["A FlowG node"]
    direction TB
    MEM["Membership layer<br/>(SWIM gossip)"]
    REP["Replication engine<br/>(broadcast + anti-entropy)"]
    ST[("Replicated stores<br/>auth · config · logs")]

    MEM -->|peer list & state| REP
    REP -->|reads / merges| ST
    ST -->|new records| REP
  end

  HTTP["HTTP management interface"]
  MEM <-->|gossip| HTTP
  REP <-->|sync batches| HTTP
  HTTP <-->|cluster traffic| PEERS["Other nodes"]
```

 - The **membership layer** discovers nodes and tracks who is alive using the
   [SWIM protocol](./consensus). It also carries the lightweight state nodes
   piggyback on each other during gossip.
 - The **replication engine** uses that membership information to
   [propagate and reconcile data](./replication) — broadcasting fresh
   control-plane writes, periodically running anti-entropy, and triggering an
   immediate incremental sync for high-volume log writes.
 - The **replicated stores** hold control-plane data as Last-Writer-Wins
   envelopes, and log entries as raw append-only records
   (see [How Data Is Replicated?](/docs/design/data-replication)).

All cluster traffic — discovery, gossip and synchronization — flows over the same
HTTP management interface, whose endpoints are described below.

## Bootstrap

When starting FlowG, it sets up a single-node cluster automatically. Every node
in the cluster needs a unique identifier. If none is given, a random one will be
generated.

Once a node joins the cluster, the cluster mesh is updated on every node of the
cluster. Replication then happens in the background via the management
interface.

Here is a diagram of the bootstrap process:

```mermaid
sequenceDiagram
  participant N1 as Node 1
  participant N2 as Node 2
  participant N3 as Node 3

  N1 ->> N1: Create single-node cluster
  N1 -->> N1: Adds "Node 1" to local cluster mesh

  critical Node 2 joins cluster via Node 1
    N2 -->> N2: Adds "Node 2" to local cluster mesh
    N2 -->> N2: Adds "Node 1" to local cluster mesh
    N2 ->> N1: Establish connection
    N1 -->> N1: Adds "Node 2" to local cluster mesh
    N1 -->> N2: Sends local cluster mesh
    N2 -->> N2: Merge remote cluster mesh with local cluster mesh
  end

  critical Node 3 joins cluster via Node 2
    N3 -->> N3: Adds "Node 3" to local cluster mesh
    N3 -->> N3: Adds "Node 2" to local cluster mesh
    N3 ->> N2: Establish connection
    N2 -->> N2: Adds "Node 3" to local cluster mesh
    N2 ->> N1: Notify new alive node "Node 3"
    N1 -->> N1: Adds "Node 3" to local cluster mesh
    N2 -->> N3: Sends local cluster mesh
    N3 -->> N3: Merge remote cluster mesh with local cluster mesh
    N3 ->> N1: Establish connection
  end
```

## Transport Endpoints

The protocol is provided on top of FlowG's HTTP management interface:

<fieldset>
<legend>Cluster Status</legend>

**Description:** Return the currently known local cluster mesh.

```
GET /cluster/nodes
```

| Response | When |
| --- | --- |
| 200 OK | On success |
| 500 Internal Server Error | On failure |

**Example:**

```json
{
  "nodes": [
    {
      "node-id": "flowg-node0",
      "endpoint": "http://<private ip address>:9113"
    },
    {
      "node-id": "flowg-node1",
      "endpoint": "http://<private ip address>:9114"
    },
    {
      "node-id": "flowg-node2",
      "endpoint": "http://<private ip address>:9115"
    }
  ]
}
```

</fieldset>

<fieldset>
<legend>Gossip (Packet mode)</legend>

**Description:** Endpoint used by the *SWIM* protocol to exchange notifications.

```
POST /cluster/gossip
Origin: <node endpoint url>
X-FlowG-ClusterKey: <optional cluster cookie>
```

| Response | When |
| --- | --- |
| 202 Accepted | On success |
| 400 Bad Request | The `Origin` header identifying the node who originated the request |
| 401 Unauthorized | The node requires authentification, and the `X-FlowG-ClusterKey` header was invalid |
| 500 Internal Server Error | On failure |

</fieldset>

<fieldset>
<legend>Gossip (stream mode)</legend>

**Description:** Bidirectional endpoint used by the *SWIM* protocol to exchange
data.

```
POST /cluster/gossip
Upgrade: flowg
X-FlowG-ClusterKey: <optional cluster cookie>
```

| Response | When |
| --- | --- |
| 101 Switching Protocols | The connection has been accepted, and the socket will be used for bi-directional exchange |
| 501 Not Implemented | The server does not support hijacking the socket, maybe there is a Reverse Proxy in between? |
| 401 Unauthorized | The node requires authentification, and the `X-FlowG-ClusterKey` header was invalid |
| 500 Internal Server Error | On failure |

</fieldset>

<fieldset>
<legend>Authentication Database synchronization</legend>

**Description:** Endpoint used to receive incremental backups from other nodes.

```
POST /cluster/sync/auth
X-FlowG-ClusterKey: <optional cluster cookie>
X-FlowG-NodeID: <remote node ID>
Transfer-Encoding: chunked
Trailer: X-FlowG-Since

... incremental backup data ...

X-FlowG-Since: <new version>
```

| Response | When |
| --- | --- |
| 401 Unauthorized | The node requires authentification, and the `X-FlowG-ClusterKey` header was invalid |
| 400 Bad Request | Missing HTTP header `X-FlowG-NodeID` or invalid data in HTTP trailer `X-FlowG-Since` |
| 500 Internal Server Error | On failure |

</fieldset>

<fieldset>
<legend>Configuration database synchronization</legend>

**Description:** Endpoint used to receive incremental backups from other nodes.

```
POST /cluster/sync/config
X-FlowG-ClusterKey: <optional cluster cookie>
X-FlowG-NodeID: <remote node ID>
Transfer-Encoding: chunked
Trailer: X-FlowG-Since

... incremental backup data ...

X-FlowG-Since: <new version>
```

| Response | When |
| --- | --- |
| 401 Unauthorized | The node requires authentification, and the `X-FlowG-ClusterKey` header was invalid |
| 400 Bad Request | Missing HTTP header `X-FlowG-NodeID` or invalid data in HTTP trailer `X-FlowG-Since` |
| 500 Internal Server Error | On failure |

</fieldset>

<fieldset>
<legend>Log database synchronization</legend>

**Description:** Endpoint used to receive incremental backups from other nodes.

```
POST /cluster/sync/log
X-FlowG-ClusterKey: <optional cluster cookie>
X-FlowG-NodeID: <remote node ID>
Transfer-Encoding: chunked
Trailer: X-FlowG-Since

... incremental backup data ...

X-FlowG-Since: <new version>
```

| Response | When |
| --- | --- |
| 401 Unauthorized | The node requires authentification, and the `X-FlowG-ClusterKey` header was invalid |
| 400 Bad Request | Missing HTTP header `X-FlowG-NodeID` or invalid data in HTTP trailer `X-FlowG-Since` |
| 500 Internal Server Error | On failure |

</fieldset>
