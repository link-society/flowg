---
sidebar_position: 2
---

# Setup

## Manually

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

Alternatively, cluster formation parameters can be set via the following
environment variables:

```bash
export FLOWG_CLUSTER_FORMATION_STRATEGY="manual"
export FLOWG_CLUSTER_FORMATION_JOIN_NODE_ID="..."
export FLOWG_CLUSTER_FORMATION_JOIN_NODE_ENDPOINT="..."
```

## Using Consul

To start a **FlowG** cluster using [Hashicorp Consul](https://developer.hashicorp.com/consul),
we need to use the `consul` cluster formation strategy:

```bash

flowg-server \
  --cluster-node-id flowg-node0 \
  --cluster-formation-strategy="consul" \
  --cluster-formation-consul-url="localhost:8500" \
  --auth-dir ./data/node0/auth \
  --log-dir ./data/node0/logs \
  --config-dir ./data/node0/config \
  --cluster-state-dir ./data/node0/state \
  --http-bind 127.0.0.1:5080 \
  --mgmt-bind 127.0.0.1:9113 \
  --syslog-bind 127.0.0.1:5514 &

flowg-server \
  --cluster-node-id flowg-node1 \
  --cluster-formation-strategy="consul" \
  --cluster-formation-consul-url="localhost:8500" \
  --auth-dir ./data/node1/auth \
  --log-dir ./data/node1/logs \
  --config-dir ./data/node1/config \
  --cluster-state-dir ./data/node1/state \
  --http-bind 127.0.0.1:5081 \
  --mgmt-bind 127.0.0.1:9114 \
  --syslog-bind 127.0.0.1:5515 &

flowg-server \
  --cluster-node-id flowg-node2 \
  --cluster-formation-strategy="consul" \
  --cluster-formation-consul-url="localhost:8500" \
  --auth-dir ./data/node2/auth \
  --log-dir ./data/node2/logs \
  --config-dir ./data/node2/config \
  --cluster-state-dir ./data/node2/state \
  --http-bind 127.0.0.1:5082 \
  --mgmt-bind 127.0.0.1:9115 \
  --syslog-bind 127.0.0.1:5516 &
```

Alternatively, automatic cluster formation can also be enabled by using the
following environment variables:

```bash
export FLOWG_CLUSTER_FORMATION_STRATEGY="consul"
export FLOWG_CLUSTER_FORMATION_CONSUL_URL="localhost:8500"
```

## Using Kubernetes (with Helm)

An Helm chart has been provided to setup **FlowG** as a `DaemonSet` alongside
[fluentd](https://www.fluentd.org/) to gather the logs of your pods:

```bash
helm install flowg ./k8s/charts/flowg -n flowg-system --create-namespace
```

> **NB:** The persistent volume used stores data on the Kubernetes node (using `hostPath`).

## Using Kubernetes (by hand)

The Kubernetes automatic cluster formation requires a `Service` so that **FlowG**
can discover all the existing nodes.

To do so, it lists (periodically) the IP addresses of the pods by querying the
proper `EndpointSlice` resource.

As such, the **FlowG** pods require a `ServiceAccount` with the following role
and role binding:

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: flowg
  namespace: flowg-system
  labels:
    app: flowg

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: flowg
  labels:
    app: flowg
rules:
  - apiGroups:
      - "discovery.k8s.io"
    resources:
      - endpointslices
    verbs:
      - list

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: flowg
  labels:
    app: flowg
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: flowg
subjects:
  - kind: ServiceAccount
    name: flowg
    namespace: flowg-system
```

Then, create the `Service` resource:

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: flowg
  labels:
    app: flowg
spec:
  type: ClusterIP
  selector:
    app: flowg
  ports:
    - name: http
      port: 5080
      targetPort: 5080
      protocol: TCP
      appProtocol: http # or 'https'
    - name: mgmt
      port: 9113
      targetPort: 9113
      protocol: TCP
      appProtocol: http # or 'https'
    - name: syslog
      port: 5514
      targetPort: 5514
      protocol: UDP
```

**FlowG** is stateful, and requires a persistent volume:

```yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: flowg-data-pvc
  labels:
    app: flowg
spec:
  storageClassName: your-storage-class
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

You can now create a `StatefulSet`:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: flowg
  labels:
    app: flowg
spec:
  selector:
    matchLabels:
      app: flowg
  template:
    metadata:
      labels:
        app: flowg
    spec:
      serviceAccountName: flowg
      containers:
        - name: flowg
          image: "linksociety/flowg:latest"
          env:
            - name: FLOWG_CLUSTER_NODE_ID
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: FLOWG_CLUSTER_FORMATION_STRATEGY
              value: "k8s"
            - name: FLOWG_CLUSTER_FORMATION_K8S_SERVICE_NAMESPACE
              value: flowg-system
            - name: FLOWG_CLUSTER_FORMATION_K8S_SERVICE_NAME
              value: flowg
            - name: FLOWG_CLUSTER_FORMATION_K8S_SERVICE_PORT_NAME
              value: mgmt
          ports:
            - containerPort: 5080
              hostPort: 5080
              protocol: TCP
            - containerPort: 9113
              hostPort: 9113
              protocol: TCP
            - containerPort: 5514
              hostPort: 5514
              protocol: UDP
          livenessProbe:
            httpGet:
              path: /health
              port: 9113
          readinessProbe:
            httpGet:
              path: /health
              port: 9113
          resources: {}
          volumeMounts:
            - name: flowg-data
              mountPath: /data
      volumes:
        - name: flowg-data
          persistentVolumeClaim:
            claimName: flowg-data-pvc
```

> **NB:** Automatic cluster formation using Kubernetes expects the node ID to be
> set to the pod's name.
