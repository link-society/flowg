---
sidebar_position: 3
---

# Helm Chart reference

This page references every available `values.yaml` option, their utility, default value(s) and behavior. 

## Base options

### `flowg`

#### `image`

<details>
<summary><code>repository</code> — <code>string</code> (Default: <code>linksociety/flowg</code>)</summary>

Repository where flowg image is sourced from. 

```yaml
flowg:
  image:
    repository: linksociety/flowg
```

</details>

<details>
<summary><code>tag</code> - <code>string</code> (Default: <code>latest</code>)</summary>

Image tag used by Docker.

```yaml
flowg:
  image:
    tag: latest
``` 

</details>

<details>
<summary><code>pullPolicy</code> - <code>string</code> (Default: <code>IfNotPresent</code>)</summary>

Image pull policy given to Kubernetes. 

```yaml
flowg:
  image:
    pullPolicy: IfNotPresent
``` 

</details>

<details>
<summary><code>pullSecrets</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Image pull secrets for private registries.

```yaml
flowg:
  image:
    pullSecrets: []
```

</details>

#### `configuration`

<details>
<summary><code>logging.verbose</code> - <code>string</code> (Default: <code>"True"</code>)</summary>

Enable or disable verbose logging.

```yaml
flowg:
  configuration:
    logging:
      verbose: "True"
```

</details>

<details>
<summary><code>logging.level</code> - <code>string</code> (Default: <code>info</code>)</summary>

Log level (e.g. debug, info, warn, error).

```yaml
flowg:
  configuration:
    logging:
      level: info
```

</details>

<details>
<summary><code>seed.secretName</code> - <code>string</code> (Default: <code>""</code>)</summary>

Kubernetes secret name containing initial/reset credentials for database seeding.

```yaml
flowg:
  configuration:
    seed:
      secretName: flowg-seed-secret
```

</details>

<details>
<summary><code>seed.initialUsernameKey</code> - <code>string</code> (Default: <code>""</code>)</summary>

Key in the secret for the initial user's username.

```yaml
flowg:
  configuration:
    seed:
      initialUsernameKey: initialUsername
```

</details>

<details>
<summary><code>seed.initialPasswordKey</code> - <code>string</code> (Default: <code>""</code>)</summary>

Key in the secret for the initial user's password.

```yaml
flowg:
  configuration:
    seed:
      initialPasswordKey: initialPassword
```

</details>

<details>
<summary><code>seed.resetUsernameKey</code> - <code>string</code> (Default: <code>""</code>)</summary>

Key in the secret for the reset user's username (applied after initial population).

```yaml
flowg:
  configuration:
    seed:
      resetUsernameKey: resetUsername
```

</details>

<details>
<summary><code>seed.resetPasswordKey</code> - <code>string</code> (Default: <code>""</code>)</summary>

Key in the secret for the reset user's password (applied after initial population).

```yaml
flowg:
  configuration:
    seed:
      resetPasswordKey: resetPassword
```

</details>

#### `podAnnotations`

<details>
<summary><code>podAnnotations</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Custom annotations to add to the FlowG pods.

```yaml
flowg:
  podAnnotations: []
```

</details>

#### `serviceAnnotations`

<details>
<summary><code>serviceAnnotations</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Custom annotations to add to the FlowG services.

```yaml
flowg:
  serviceAnnotations: []
```

</details>

#### `podSecurityContext`

<details>
<summary><code>runAsNonRoot</code> - <code>boolean</code> (Default: <code>true</code>)</summary>

Requires the pod to run as a non-root user.

```yaml
flowg:
  podSecurityContext:
    runAsNonRoot: true
```

</details>

<details>
<summary><code>runAsUser</code> - <code>integer</code> (Default: <code>10001</code>)</summary>

The UID to run the container process as.

```yaml
flowg:
  podSecurityContext:
    runAsUser: 10001
```

</details>

<details>
<summary><code>runAsGroup</code> - <code>integer</code> (Default: <code>10001</code>)</summary>

The GID to run the container process as.

```yaml
flowg:
  podSecurityContext:
    runAsGroup: 10001
```

</details>

<details>
<summary><code>fsGroup</code> - <code>integer</code> (Default: <code>10001</code>)</summary>

The filesystem group ID for volumes.

```yaml
flowg:
  podSecurityContext:
    fsGroup: 10001
```

</details>

#### `securityContext`

<details>
<summary><code>securityContext</code> - <code>object</code> (Default: <code>{}</code>)</summary>

Container-level security context (e.g. allowPrivilegeEscalation, capabilities).

```yaml
flowg:
  securityContext: {}
```

</details>

#### `probe`

<details>
<summary><code>livenessProbe</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Custom configuration for liveness probe settings (initialDelaySeconds, timeoutSeconds, etc.).

```yaml
flowg:
  probe:
    livenessProbe: []
```

</details>

<details>
<summary><code>readinessProbe</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Custom configuration for readiness probe settings (initialDelaySeconds, timeoutSeconds, etc.).

```yaml
flowg:
  probe:
    readinessProbe: []
```

</details>

#### `nodeSelectorTolerations`

<details>
<summary><code>nodeSelectorTolerations</code> - <code>array</code> (Default: <code>null</code>)</summary>

Tolerations for pod scheduling on tainted nodes.

```yaml
flowg:
  nodeSelectorTolerations: []
```

</details>

#### `storage`

<details>
<summary><code>size</code> - <code>string</code> (Default: <code>10Gi</code>)</summary>

Storage size requested for the data volume PVC.

```yaml
flowg:
  storage:
    size: 10Gi
```

</details>

<details>
<summary><code>storageClassName</code> - <code>string</code> (Default: <code>""</code>)</summary>

Storage class name for provisioning the PVC.

```yaml
flowg:
  storage:
    storageClassName: ""
```

</details>

#### `httpPort`

<details>
<summary><code>httpPort</code> - <code>integer</code> (Default: <code>5080</code>)</summary>

Port on which the HTTP server listens inside the container.

```yaml
flowg:
  httpPort: 5080
```

</details>

#### `syslog`

<details>
<summary><code>port</code> - <code>integer</code> (Default: <code>5514</code>)</summary>

Port on which the syslog server listens.

```yaml
flowg:
  syslog:
    port: 5514
```

</details>

<details>
<summary><code>protocol</code> - <code>string</code> (Default: <code>UDP</code>)</summary>

Protocol used by the syslog server (e.g. UDP, TCP).

```yaml
flowg:
  syslog:
    protocol: UDP
```

</details>

#### `resources`

<details>
<summary><code>resources</code> - <code>object</code> (Default: <code>{}</code>)</summary>

CPU and memory resource requests and limits for the FlowG container.

```yaml
flowg:
  resources: {}
```

</details>

#### `serviceAccount`

<details>
<summary><code>create</code> - <code>boolean</code> (Default: <code>true</code>)</summary>

Specifies whether a service account should be created.

```yaml
flowg:
  serviceAccount:
    create: true
```

</details>

<details>
<summary><code>annotations</code> - <code>object</code> (Default: <code>{}</code>)</summary>

Annotations to add to the service account.

```yaml
flowg:
  serviceAccount:
    annotations: {}
```

</details>

<details>
<summary><code>name</code> - <code>string</code> (Default: <code>""</code>)</summary>

The name of the service account to use.

```yaml
flowg:
  serviceAccount:
    name: ""
```

</details>

#### `service`

<details>
<summary><code>type</code> - <code>string</code> (Default: <code>LoadBalancer</code>)</summary>

Kubernetes service type (e.g. LoadBalancer, NodePort, ClusterIP).

```yaml
flowg:
  service:
    type: LoadBalancer
```

</details>

<details>
<summary><code>nodePort.httpPort</code> - <code>integer</code> (Default: <code>30080</code>)</summary>

NodePort assigned for HTTP traffic when service type is NodePort/LoadBalancer.

```yaml
flowg:
  service:
    nodePort:
      httpPort: 30080
```

</details>

<details>
<summary><code>nodePort.syslogPort</code> - <code>integer</code> (Default: <code>32514</code>)</summary>

NodePort assigned for syslog traffic when service type is NodePort/LoadBalancer.

```yaml
flowg:
  service:
    nodePort:
      syslogPort: 32514
```

</details>


