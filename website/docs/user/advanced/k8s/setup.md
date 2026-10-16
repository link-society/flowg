---
sidebar_position: 3
---

# Setting up a Kubernetes cluster

This guide helps you setup a flowG Kubernetes cluster using our **Helm charts**. 

This guide splits into 3 different branches, following the needs of your cluster.

1. [Single-node flowg instance with BadgerDB as a backend](), if you only need a small installation. 
2. [Multi-node cluster using foundationDB](), if you need multiple instances running on different nodes.
3. [Multi-node cluster with fluentD](), if you dont want to bring your own log routing system. 

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