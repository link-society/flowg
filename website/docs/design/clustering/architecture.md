---
sidebar_position: 2
---

# Architecture

## Bootstrap

When starting FlowG, it sets up a single-node cluster automatically. Every node
in the cluster needs a unique identifier. If none is given, a random one will be
generated.

You can tell your instance to join an existing cluster by indicating the name of
the node to join, and its management endpoint.

> **NB:** Automatic cluster formation is planned via:
>
> - Kubernetes headless `Service` resource
> - Consul's service mesh
> - DNS discovery

Here is an example of how to start a 3-nodes cluster:

```bash
flowg-server \
  --cluster-node-id flowg-node0 \
  --cluster-formation-strategy="manual" \
  --auth-dir ./data/node0/auth \
  --log-dir ./data/node0/logs \
  --config-dir ./data/node0/config \
  --cluster-state-dir ./data/node0/state \
  --http-bind 127.0.0.1:5080 \
  --mgmt-bind 127.0.0.1:9113 \
  --syslog-bind 127.0.0.1:5514 &

flowg-server \
  --cluster-node-id flowg-node1 \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node0 \
  --cluster-formation-manual-join-endpoint http://localhost:9113 \
  --auth-dir ./data/node1/auth \
  --log-dir ./data/node1/logs \
  --config-dir ./data/node1/config \
  --cluster-state-dir ./data/node1/state \
  --http-bind 127.0.0.1:5081 \
  --mgmt-bind 127.0.0.1:9114 \
  --syslog-bind 127.0.0.1:5515 &

flowg-server \
  --cluster-node-id flowg-node2 \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node1 \
  --cluster-formation-manual-join-endpoint http://localhost:9114 \
  --auth-dir ./data/node2/auth \
  --log-dir ./data/node2/logs \
  --config-dir ./data/node2/config \
  --cluster-state-dir ./data/node2/state \
  --http-bind 127.0.0.1:5082 \
  --mgmt-bind 127.0.0.1:9115 \
  --syslog-bind 127.0.0.1:5516 &
```

> **NB:** Don't use `&` to run FlowG in the background, this is just an example.

You can also enable authentication between nodes by using a secret key that
each node requires:

```bash
cookie=$(openssl rand -hex 32)

flowg-server \
  --cluster-node-id flowg-node0 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --auth-dir ./data/node0/auth \
  --log-dir ./data/node0/logs \
  --config-dir ./data/node0/config \
  --cluster-state-dir ./data/node0/state \
  --http-bind 127.0.0.1:5080 \
  --mgmt-bind 127.0.0.1:9113 \
  --syslog-bind 127.0.0.1:5514 &

flowg-server \
  --cluster-node-id flowg-node1 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node0 \
  --cluster-formation-manual-join-endpoint http://localhost:9113 \
  --auth-dir ./data/node1/auth \
  --log-dir ./data/node1/logs \
  --config-dir ./data/node1/config \
  --cluster-state-dir ./data/node1/state \
  --http-bind 127.0.0.1:5081 \
  --mgmt-bind 127.0.0.1:9114 \
  --syslog-bind 127.0.0.1:5515 &

flowg-server \
  --cluster-node-id flowg-node2 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node1 \
  --cluster-formation-manual-join-endpoint http://localhost:9114 \
  --auth-dir ./data/node2/auth \
  --log-dir ./data/node2/logs \
  --config-dir ./data/node2/config \
  --cluster-state-dir ./data/node2/state \
  --http-bind 127.0.0.1:5082 \
  --mgmt-bind 127.0.0.1:9115 \
  --syslog-bind 127.0.0.1:5516 &
```

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
