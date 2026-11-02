---
sidebar_position: 7
---

# Tuning

This page explains how to size FlowG for a given log throughput.

## Workers

FlowG processes data concurrently with fixed pools of workers. The `workers`
section of the [configuration file](/docs/technical/configuration) sets the
size of each pool:

```hcl
workers {
  # env: FLOWG_WORKERS_PIPELINES (default: 64)
  pipelines = 64

  # env: FLOWG_WORKERS_SYSLOG (default: 64)
  syslog = 64

  # env: FLOWG_WORKERS_BADGERDB (default: 64)
  badgerdb = 64
}
```

 - **pipelines**: number of records processed by the pipelines at the same
   time, whatever their input (HTTP API or Syslog)
 - **syslog**: number of received Syslog messages processed at the same time
 - **badgerdb**: number of concurrent operations per BadgerDB database (the
   auth, config and log databases each have their own pool), ignored by the
   FoundationDB backend

Every count must be at least 1.

With the Helm charts, the same settings are available under
`flowg.configuration.workers` (`badgerdb` only in the single-node chart):

```yaml
flowg:
  configuration:
    workers:
      pipelines: "64"
      syslog: "64"
      badgerdb: "64"
```

## Sizing

The following table gives the smallest container that sustains a given
throughput, measured with the FlowG stress test suite:

| Input Protocol | Maximum Throughput | Minimum CPU requirements | Minimum RAM requirements |
| --- | --- | --- | --- |
| HTTP | 2k logs/s | 0.128 CPU | 256 MiB |
| HTTP | 5k logs/s | 0.256 CPU | 512 MiB |
| HTTP | 10k logs/s | 0.512 CPU | 1 GiB |
| Syslog UDP | 2k logs/s | 0.128 CPU | 256 MiB |
| Syslog UDP | 5k logs/s | 0.256 CPU | 512 MiB |
| Syslog UDP | 10k logs/s | 1 CPU | 2 GiB |
| Syslog TCP | 2k logs/s | 0.128 CPU | 256 MiB |
| Syslog TCP | 5k logs/s | 0.256 CPU | 512 MiB |
| Syslog TCP | 10k logs/s | 0.512 CPU | 1 GiB |

Heavier pipelines (more transformers, forwarders, etc...) need more resources.

:::warning Syslog over UDP

UDP has no flow control: when FlowG cannot keep up, the messages it does not
have the time to read are dropped by the operating system, without any error.
Syslog over UDP therefore needs more CPU than Syslog over TCP for the same
throughput. Prefer TCP when you cannot afford to lose logs.

:::

## BadgerDB

For the BadgerDB storage backend, use the following memtable settings for each
container size, based on the tuned profiles in the FlowG stress test suite:

| CPU | RAM | Memtable size (MB) | Memtables count |
| --- | --- | --- | --- |
| 0.128 CPU | 256 MiB | 8 | 5 |
| 0.256 CPU | 512 MiB | 16 | 5 |
| 0.512 CPU | 1 GiB | 32 | 5 |
| 1 CPU | 2 GiB | 64 | 5 |

These values apply to **each** of the auth, config and log databases. The stress
tests tune the memtable size for all three databases and keep the default count
of 5. The default memtable size is 64 MB; reduce it for smaller containers as
shown above. Memory used by memtables is roughly their size times their count,
per database, in addition to the memory needed by the rest of FlowG.

In the [configuration file](/docs/technical/configuration), set
`auth_memtable_size`, `config_memtable_size` and `log_memtable_size` under
`storage.backend "badgerdb"`. The corresponding count settings are
`auth_memtables_count`, `config_memtables_count` and `log_memtables_count`.
With the [single-node Helm chart](/docs/technical/helm), use
`flowg.configuration.badger` with `authMemtableSize`, `configMemtableSize`,
`logMemtableSize`, `authMemtablesCount`, `configMemtablesCount` and
`logMemtablesCount`.

## Memory limit

By default, the Go runtime does not know about the memory limit of the
container: under load, its garbage collector lets the memory grow up to twice
what FlowG actually uses, and the container gets killed.

Set the `GOMEMLIMIT` environment variable to about 80% of the container memory
limit, the remaining 20% being used by memory the Go runtime does not manage
(the VRL transformers and the files mapped by BadgerDB).

With Docker:

```bash
docker run \
  --cpus 0.256 \
  --memory 512m \
  -e GOMEMLIMIT=400MiB \
  -p 5080:5080/tcp \
  -p 9113:9113/tcp \
  -p 5514:5514/udp \
  -v flowg-data:/data \
  linksociety/flowg:latest
```

With the Helm charts, set `flowg.goMemLimit` along with the memory limit in
`flowg.resources`:

```yaml
flowg:
  resources:
    limits:
      cpu: 256m
      memory: 512Mi
  goMemLimit: "400MiB"
```
