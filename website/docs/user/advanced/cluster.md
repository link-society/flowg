---
sidebar_position: 2
---

# Setting up a Cluster

By default, **FlowG** uses [BadgerDB](https://github.com/dgraph-io/badger), an
embedded key/value store, which limits a deployment to a single node.

To run **FlowG** in a highly-available, horizontally scalable cluster, you can
switch the storage backend to [FoundationDB](https://www.foundationdb.org/), a
distributed key/value store. Every **FlowG** node then becomes stateless and
reads/writes the same shared FoundationDB cluster, so any node can ingest and
query logs.

> **NB:** All nodes in the cluster **must** share the same `FLOWG_SECRET_KEY`.
> This key is used to sign the JSON Web Tokens handed out on login. If the nodes
> use different keys, a token issued by one node will be rejected by the others.

> **NB:** **FlowG** is built against the FoundationDB **7.3** client library
> (bundled in the release archives and in the Docker image). Make sure your
> FoundationDB **server** uses a compatible **7.3.x** version.

## Without Docker

### 1. Set up a FoundationDB cluster

Follow the [official FoundationDB documentation](https://apple.github.io/foundationdb/administration.html)
to deploy a FoundationDB cluster. Once configured, FoundationDB generates a
**cluster file** (usually located at `/etc/foundationdb/fdb.cluster`) that
describes how to reach the coordinators.

### 2. Install the FoundationDB client library on each FlowG host

**FlowG** connects to FoundationDB through the client library, which is bundled
in the release archives under `lib/`. Install it into a directory searched by
the dynamic linker.

On **Linux** (`libfdb_c.so`):

```bash
sudo install -m 755 ./lib/libfdb_c.so /usr/local/lib/libfdb_c.so
sudo ldconfig
```

On **macOS** (`libfdb_c.dylib`):

```bash
sudo install -m 755 ./lib/libfdb_c.dylib /usr/local/lib/libfdb_c.dylib
```

Then copy the cluster file from your FoundationDB cluster to each **FlowG** host,
for example at `/etc/foundationdb/fdb.cluster`.

### 3. Configure each FlowG node

On every node, create a configuration file `/etc/flowg/config.hcl` that selects
the FoundationDB backend:

```hcl
services {
  http {
    bind = "0.0.0.0:5080"
  }

  management {
    bind = "0.0.0.0:9113"
  }

  syslog {
    bind = "0.0.0.0:5514"
  }
}

storage {
  backend "foundationdb" {
    cluster_file = "/etc/foundationdb/fdb.cluster"
    ## Or:
    # connection_string = "..."
    key_space    = "flowg"
  }
}
```

> **NB:** The `key_space` is optional and defaults to `flowg`. It lets multiple
> **FlowG** clusters share the same FoundationDB cluster without interfering
> with each other.

### 4. Start each node

On every node, export the **shared** secret key and start the server:

```bash
export FLOWG_SECRET_KEY="a-long-random-shared-secret"
flowg-server --config /etc/flowg/config.hcl
```

Every node is now part of the same cluster. You can place them behind a load
balancer to distribute ingestion and query traffic.

## With Docker (Docker Compose)

The following `docker-compose.yml` starts a single-node FoundationDB, initializes
it, and runs a 3-node **FlowG** cluster on top of it:

```yaml
name: flowg-cluster

services:
  foundationdb:
    image: foundationdb/foundationdb:7.3.77
    environment:
      FDB_NETWORKING_MODE: "container"
      FDB_PORT: "4500"
      FDB_CLUSTER_FILE: "/etc/foundationdb/fdb.cluster"
    volumes:
      - fdb-data:/var/fdb/data
      - fdb-config:/etc/foundationdb

  foundationdb-init:
    image: foundationdb/foundationdb:7.3.77
    depends_on:
      foundationdb:
        condition: service_started
    restart: "no"
    entrypoint: ["/bin/bash", "-ec"]
    command:
      - |
        until test -s /etc/foundationdb/fdb.cluster; do
          sleep 1
        done

        until fdbcli -C /etc/foundationdb/fdb.cluster --exec "status" --timeout 2; do
          fdbcli -C /etc/foundationdb/fdb.cluster --exec "configure new single ssd" --timeout 10 || true
        done

        fdbcli -C /etc/foundationdb/fdb.cluster --exec "status"
    volumes:
      - fdb-config:/etc/foundationdb

  flowg-node-1:
    image: linksociety/flowg:latest
    depends_on:
      foundationdb-init:
        condition: service_completed_successfully
    environment:
      FLOWG_SECRET_KEY: "a-long-random-shared-secret"
      FLOWG_STORAGE_BACKEND: "foundationdb"
      FLOWG_FOUNDATIONDB_CLUSTER_FILE: "/etc/foundationdb/fdb.cluster"
      FLOWG_FOUNDATIONDB_KEY_SPACE: "flowg"
    ports:
      - "5080:5080/tcp"
      - "5514:5514/udp"
    volumes:
      - fdb-config:/etc/foundationdb:ro

  flowg-node-2:
    image: linksociety/flowg:latest
    depends_on:
      foundationdb-init:
        condition: service_completed_successfully
    environment:
      FLOWG_SECRET_KEY: "a-long-random-shared-secret"
      FLOWG_STORAGE_BACKEND: "foundationdb"
      FLOWG_FOUNDATIONDB_CLUSTER_FILE: "/etc/foundationdb/fdb.cluster"
      FLOWG_FOUNDATIONDB_KEY_SPACE: "flowg"
    ports:
      - "5081:5080/tcp"
      - "5515:5514/udp"
    volumes:
      - fdb-config:/etc/foundationdb:ro

  flowg-node-3:
    image: linksociety/flowg:latest
    depends_on:
      foundationdb-init:
        condition: service_completed_successfully
    environment:
      FLOWG_SECRET_KEY: "a-long-random-shared-secret"
      FLOWG_STORAGE_BACKEND: "foundationdb"
      FLOWG_FOUNDATIONDB_CLUSTER_FILE: "/etc/foundationdb/fdb.cluster"
      FLOWG_FOUNDATIONDB_KEY_SPACE: "flowg"
    ports:
      - "5082:5080/tcp"
      - "5516:5514/udp"
    volumes:
      - fdb-config:/etc/foundationdb:ro

volumes:
  fdb-data:
  fdb-config:
```

A few things to note about this configuration:

 - `FDB_NETWORKING_MODE` is set to `container` so FoundationDB advertises its
   container address in the cluster file, making it reachable by the other
   containers on the same Docker network.
 - The `fdb-config` volume holds the cluster file generated by FoundationDB. It
   is shared with the **FlowG** nodes (mounted read-only) so they can connect to
   the cluster.
 - `foundationdb-init` waits for FoundationDB to be up, then configures a fresh
   database. The **FlowG** nodes only start once this initialization has
   completed successfully.
 - Every **FlowG** node shares the same `FLOWG_SECRET_KEY`.

Start the cluster with:

```bash
docker compose up
```

The three nodes are now available on ports `5080`, `5081` and `5082`.

> **NB:** The `configure new single ssd` command creates a **single-node**
> database with no redundancy, which is only suitable for testing. For a
> production cluster, deploy multiple FoundationDB processes and use a redundant
> mode such as `double` or `triple`. See the
> [FoundationDB documentation](https://apple.github.io/foundationdb/configuration.html)
> for details.
