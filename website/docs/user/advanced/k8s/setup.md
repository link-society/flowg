---
sidebar_position: 3
---

# Setting up a Kubernetes cluster

This guide helps you setup a flowG Kubernetes cluster using our **Helm charts**. 

This guide splits into 3 different branches, following the needs of your cluster.

1. [Single-node flowg instance with BadgerDB as a backend](), if you only need a small installation. 
2. [Multi-node cluster using foundationDB]() (and optionnally an observability stack), if you need multiple instances running 

## Prerequisites 

- A working **Kubernetes cluster** 
- Helm 

## 1. Single-node flowg instance 

First, clone the repository containing our helm charts:

```bash 
git clone https://github.com/link-society/flowg
cd flowg/k8s/charts
```

> **Configuration**  
> The configuration can be found inside `values.yaml`. The default configuration normally applies to all needs, but you can still customize it if you want. See [configuration](configuration.md)


To install the chart, simply run:
```bash
helm install flowg flowg-singlenode -n flowg-system --create-namespace
```

This will create these resources into the chosen namespace:
- (If HTTPS is enabled on web or MGMT endpoint) 1 or 2 cert-manager certificates
- A `Service` used to access flowg's web interface
- A `Service` used to perform healthchecks on flowg's mgmt endpoint
- A `ConfigMap` containing the applied configuration as environment variables, sourced from `values.yaml` 
- A `ServiceAccount` for flowG
- A `StatefulSet` with flowG's pod and 1 PVC for its storage

## 2. Multi-node cluster, with/out Fluentd 

### Installing FoundationDB CRD

In order to make flowG work with FoundationDB, you first need to intall the [FoundationDB Kubernetes Operator](https://github.com/foundationdb/fdb-kubernetes-operator). 

To do so, a simple way is to use the Helm chart they provide.

```bash
# Clone the repository 
git clone https://github.com/foundationdb/fdb-kubernetes-operator
cd fdb-kubernetes-operator

# Install FoundationDB CRD
# We install it in the same namespace as flowG 
helm install foundationdb charts/fdb-opeator -n flowg-system --create-namespace
``` 

### Setting up flowG

The instructions are close to the ones for the other charts. First, you need to clone the repository:
```bash
git clone https://github.com/link-society/flowg
cd flowg
```

> **Configuration**  
> The configuration can be found inside `values.yaml`. The default configuration normally applies to all needs, but you can still customize it if you want. See [configuration](configuration.md)

And then install the chart following your needs:

- **If you don't need fluentD**: 
```bash
helm install flowg flowg-multinode -n flowg-system --create-namespace
```

- **If you want an observability stack based on fluentd**: 
```bash
helm install flowg flowg-observability -n flowg-system --create-namespace
``` 

> **Fluentd**  
> Fluentd manages all the log routing to flowG by gathering ALL the logfiles specified in the `logdir` in `values.yaml`, except fluentd's own logs, and sends it to flowG's syslog server. You can learn more [here]() about how to configure it.


> **Important note**  
> The cluster will take some time to get brought up. This is normal. FoundationDB takes time reconciling their cluster and creating the cluster file. An init container waiting for it will be waiting for it to exist before running flowG (to avoid it crashing on startup).

This chart will set these resources up inside your namespace:
- (If HTTPS is enabled on web or MGMT endpoint) 1 or 2 cert-manager certificates
- A `Service` used to access flowg's web interface
- A `Service` used to perform healthchecks on flowg's mgmt endpoint
- A `ConfigMap` containing the applied configuration as environment variables, sourced from `values.yaml` 
- A `ServiceAccount` for flowG
- A `Deployment` with flowG's pods
- A `FoundationDBCluster` for flowG's storage, with its own RBAC setup
- A `Secret` containing a random-generated key in order for flowG instances to share the same JWT secret key.
