---
sidebar_position: 2
---

# Setting up HTTPS

There are two main ways of setting up HTTPS in flowg. The first lets you use any provider you want, and only needs you to store the certificate and private key inside a Kubernetes secret. The second uses [cert-manager](https://cert-manager.io/docs) to auto generate and manage the TLS certificate.

:::note  
These exact instructions will only work for the web enpoint (port 5080 by default). If you need to setup TLS encryption for the MGMT endpoint (port 9113), just replace `https` by `https.mgmt` in the `values.yaml` setup. You can have both enabled.
:::
## 1. Using cert-manager

You will first need to create an `Issuer` and `ClusterIssuer` resource. To setup cert-manager and create these resources, see [the official Cert-Manager documentation](https://cert-manager.io/docs/getting-started/)

You will need to redeploy the chart with these values:

```yaml

# ... 
https:
    enabled: true

    certificateFrom:
      ## Manually issued certificate stored in a `Secret`:
      # secretRef:
      #   ## Expected keys: `tls.crt` & `tls.key`
      #   name: logs-example-com-tls

      ## Automatically issued with Cert-Manager
      certmanager:
        commonName: logs.example.com
        issuerRef:
          name: ca-issuer      # adjust to your setup
          kind: ClusterIssuer  # adjust to your setup
```

## 2. Using a custom secret

If you want instead to use an external certificate management system, or even manage them manually, the Helm chart allows you to connect it to a secret.

In order to set that up, you will first need to create a `Secret` in the same namespace as flowg.  
This secret will need to have two keys:
 - `tls.crt`, containing the certificate,
 - `tls.key`, containing the private key for the certificate. 

This is important for the keys to match their exact names, insted flowg will not be able to discover them.

Once the `Secret` created, just redeploy the Helm chart with this `values.yaml` setup: 

```yaml 

# ...

https:
    enabled: true
    
    certificateFrom: 
        secretRef:
            name: my-secret # Change according to your setup
``` 

Once this is done, everything should be setup and ready!

