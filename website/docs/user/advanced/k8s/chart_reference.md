---
sidebar_position: 3
toc_min_heading_level: 2
toc_max_heading_level: 2
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Charts configuration reference

This page references every available `values.yaml` option, their utility, default value(s) and behavior. 

<Tabs>
  <TabItem value="base" label="Base options" default>


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

#### Multi-node flowg options

<details>
<summary><code>replicas</code> - <code>integer</code> (Default: <code>4</code>)</summary>

Number of FlowG replicas in the multi-node deployment.

```yaml
flowg:
  replicas: 4
```

</details>

<details>
<summary><code>secretKey.useExistingSecret</code> - <code>boolean</code> (Default: <code>false</code>)</summary>

Whether to use an existing Kubernetes secret for `FLOWG_SECRET_KEY`.

```yaml
flowg:
  secretKey:
    useExistingSecret: false
```

</details>

<details>
<summary><code>secretKey.secretRef</code> - <code>string</code> (Default: <code>flowg-shared-key</code>)</summary>

Name of the secret/key used for shared session encryption across replicas.

```yaml
flowg:
  secretKey:
    secretRef: flowg-shared-key
```

</details>

  </TabItem>
  <TabItem value="https" label="HTTPS">

The main `https` key enables TLS encryption on the web endpoint (port 5080 by default), while `https.mgmt` enables it on the MGMT endpoint (port 9113).

### `https` / `https.mgmt`

<details>
<summary><code>enabled</code> - <code>boolean</code> (Default: <code>false</code>)</summary>

Enables HTTPS encryption on the web endpoint (`https`) or management endpoint (`https.mgmt`).

```yaml
flowg:
  https:
    enabled: true
```

</details>

<details>
<summary><code>certificateFrom.secretRef.name</code> - <code>string</code> (Default: <code>""</code>)</summary>

Name of an existing Kubernetes Secret containing `tls.crt` and `tls.key` for custom certificate management.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      secretRef:
        name: flowg-tls
```

</details>

<details>
<summary><code>certificateFrom.certmanager.privateKey.algorithm</code> - <code>string</code> (Default: <code>RSA</code>)</summary>

Private key algorithm generated by Cert-Manager (e.g. RSA, ECDSA).

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        privateKey:
          algorithm: RSA
```

</details>

<details>
<summary><code>certificateFrom.certmanager.privateKey.encoding</code> - <code>string</code> (Default: <code>PKCS1</code>)</summary>

Private key encoding generated by Cert-Manager.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        privateKey:
          encoding: PKCS1
```

</details>

<details>
<summary><code>certificateFrom.certmanager.privateKey.size</code> - <code>integer</code> (Default: <code>2048</code>)</summary>

Private key size in bits generated by Cert-Manager.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        privateKey:
          size: 2048
```

</details>

<details>
<summary><code>certificateFrom.certmanager.dnsNames</code> - <code>array</code> (Default: <code>[]</code>)</summary>

DNS SANs for the automatically generated certificate.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        dnsNames:
          - logs.example.com
```

</details>

<details>
<summary><code>certificateFrom.certmanager.ipAddresses</code> - <code>array</code> (Default: <code>["127.0.0.1"]</code>)</summary>

IP SANs for the automatically generated certificate.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        ipAddresses:
          - "127.0.0.1"
```

</details>

<details>
<summary><code>certificateFrom.certmanager.issuerRef.name</code> - <code>string</code> (Default: <code>selfsigned-issuer</code>)</summary>

Name of the Cert-Manager Issuer or ClusterIssuer.

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        issuerRef:
          name: selfsigned-issuer
          kind: Issuer
```

</details>

<details>
<summary><code>certificateFrom.certmanager.issuerRef.kind</code> - <code>string</code> (Default: <code>Issuer</code>)</summary>

Kind of the Cert-Manager issuer (`Issuer` or `ClusterIssuer`).

```yaml
flowg:
  https:
    enabled: true
    certificateFrom:
      certmanager:
        issuerRef:
          name: selfsigned-issuer
          kind: Issuer
```

</details>

  </TabItem>
  <TabItem value="foundationdb" label="FoundationDB">


These options configure the FoundationDB cluster used by the multi-node FlowG deployment.

### `foundationdb`

<details>
<summary><code>processCounts.storage</code> - <code>integer</code> (Default: <code>0</code>)</summary>

Number of FoundationDB storage processes per FDB process group (`-1` is none, `0` lets FDB decide).

```yaml
foundationdb:
  processCounts:
    storage: 0
```

</details>

<details>
<summary><code>processCounts.log</code> - <code>integer</code> (Default: <code>-1</code>)</summary>

Number of FoundationDB log processes per FDB process group.

```yaml
foundationdb:
  processCounts:
    log: -1
```

</details>

<details>
<summary><code>processCounts.stateless</code> - <code>integer</code> (Default: <code>-1</code>)</summary>

Number of FoundationDB stateless processes per FDB process group.

```yaml
foundationdb:
  processCounts:
    stateless: -1
```

</details>

<details>
<summary><code>redundancyMode</code> - <code>string</code> (Default: <code>single</code>)</summary>

FDB redundancy mode (e.g. single, double, triple).

```yaml
foundationdb:
  redundancyMode: single
```

</details>

<details>
<summary><code>resourceRequests.worker.cpu</code> - <code>string/number</code> (Default: <code>1</code>)</summary>

CPU resource request for FDB worker processes.

```yaml
foundationdb:
  resourceRequests:
    worker:
      cpu: 1
```

</details>

<details>
<summary><code>resourceRequests.worker.memory</code> - <code>string</code> (Default: <code>128Mi</code>)</summary>

Memory resource request for FDB worker processes.

```yaml
foundationdb:
  resourceRequests:
    worker:
      memory: 128Mi
```

</details>

<details>
<summary><code>resourceRequests.sidecar.cpu</code> - <code>string/number</code> (Default: <code>1</code>)</summary>

CPU resource request for FDB sidecar container.

```yaml
foundationdb:
  resourceRequests:
    sidecar:
      cpu: 1
```

</details>

<details>
<summary><code>resourceRequests.sidecar.memory</code> - <code>string</code> (Default: <code>128Mi</code>)</summary>

Memory resource request for FDB sidecar container.

```yaml
foundationdb:
  resourceRequests:
    sidecar:
      memory: 128Mi
```

</details>

<details>
<summary><code>resourceLimits.worker.cpu</code> - <code>string/number</code> (Default: <code>1</code>)</summary>

CPU resource limit for FDB worker processes.

```yaml
foundationdb:
  resourceLimits:
    worker:
      cpu: 1
```

</details>

<details>
<summary><code>resourceLimits.worker.memory</code> - <code>string</code> (Default: <code>128Mi</code>)</summary>

Memory resource limit for FDB worker processes.

```yaml
foundationdb:
  resourceLimits:
    worker:
      memory: 128Mi
```

</details>

<details>
<summary><code>resourceLimits.sidecar.cpu</code> - <code>string/number</code> (Default: <code>1</code>)</summary>

CPU resource limit for FDB sidecar container.

```yaml
foundationdb:
  resourceLimits:
    sidecar:
      cpu: 1
```

</details>

<details>
<summary><code>resourceLimits.sidecar.memory</code> - <code>string</code> (Default: <code>128Mi</code>)</summary>

Memory resource limit for FDB sidecar container.

```yaml
foundationdb:
  resourceLimits:
    sidecar:
      memory: 128Mi
```

</details>

<details>
<summary><code>sidecarContainer.livenessProbe</code> - <code>boolean</code> (Default: <code>true</code>)</summary>

Enable or disable liveness probe for the FDB sidecar container.

```yaml
foundationdb:
  sidecarContainer:
    livenessProbe: true
```

</details>

<details>
<summary><code>sidecarContainer.readinessProbe</code> - <code>boolean</code> (Default: <code>false</code>)</summary>

Enable or disable readiness probe for the FDB sidecar container.

```yaml
foundationdb:
  sidecarContainer:
    readinessProbe: false
```

</details>

<details>
<summary><code>storageRequest</code> - <code>string</code> (Default: <code>16Gi</code>)</summary>

Storage size requested for each FoundationDB PV.

```yaml
foundationdb:
  storageRequest: 16Gi
```

</details>

  </TabItem>
  <TabItem value="fluentd" label="Fluentd">

These options configure the Fluentd daemonset used in the observability stack chart (`flowg-observability`).

### `fluentd`

<details>
<summary><code>image.repository</code> - <code>string</code> (Default: <code>fluent/fluentd-kubernetes-daemonset</code>)</summary>

Repository for the Fluentd image.

```yaml
fluentd:
  image:
    repository: fluent/fluentd-kubernetes-daemonset
```

</details>

<details>
<summary><code>image.tag</code> - <code>string</code> (Default: <code>v1-debian-syslog</code>)</summary>

Image tag for Fluentd.

```yaml
fluentd:
  image:
    tag: v1-debian-syslog
```

</details>

<details>
<summary><code>image.pullPolicy</code> - <code>string</code> (Default: <code>IfNotPresent</code>)</summary>

Image pull policy for Fluentd.

```yaml
fluentd:
  image:
    pullPolicy: IfNotPresent
```

</details>

<details>
<summary><code>image.pullSecrets</code> - <code>array</code> (Default: <code>[]</code>)</summary>

Image pull secrets for private container registries.

```yaml
fluentd:
  image:
    pullSecrets: []
```

</details>

<details>
<summary><code>logdir</code> - <code>string</code> (Default: <code>/var/log/pods</code>)</summary>

Directory on nodes where pod logs are gathered from.

```yaml
fluentd:
  logdir: /var/log/pods
```

</details>

<details>
<summary><code>resources</code> - <code>object</code> (Default: <code>{}</code>)</summary>

Resource requests and limits for Fluentd pods.

```yaml
fluentd:
  resources: {}
```

</details>

<details>
<summary><code>serviceAccount.create</code> - <code>boolean</code> (Default: <code>true</code>)</summary>

Specifies whether a service account should be created for Fluentd.

```yaml
fluentd:
  serviceAccount:
    create: true
```

</details>

<details>
<summary><code>serviceAccount.annotations</code> - <code>object</code> (Default: <code>{}</code>)</summary>

Annotations to add to Fluentd's service account.

```yaml
fluentd:
  serviceAccount:
    annotations: {}
```

</details>

<details>
<summary><code>serviceAccount.name</code> - <code>string</code> (Default: <code>""</code>)</summary>

Name of the Fluentd service account.

```yaml
fluentd:
  serviceAccount:
    name: ""
```

</details>

  </TabItem>
</Tabs>
