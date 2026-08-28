---
sidebar_position: 11
---

# Azure Monitor

This forwarder is used to send a log record to
[Azure Monitor Logs](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/data-platform-logs),
via the
[Logs Ingestion API](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/logs-ingestion-api-overview).

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str endpoint "URL of the Data Collection Endpoint to upload the records to"
    str token[1] "Microsoft Entra ID access token used to authenticate"
    str expires_on "Expiry date of the token, in RFC3339 format"
    str rule_id "Immutable ID of the Data Collection Rule"
    str stream_name "Name of the input stream declared in the Data Collection Rule"
    bool allow_insecure "Skip TLS certificate verification"
  }
```

:::note

1. The token is **NOT** encrypted in the database. FlowG does not renew it
   either: the forwarder must be reconfigured with a fresh token before it
   expires.
2. The Data Collection Rule, its input stream and the destination table must
   exist prior to sending records, they are not created by the forwarder.
3. `allow_insecure` disables TLS certificate verification, it is meant to
   target a local emulator or a test server, it should not be used against
   the Azure Monitor API.

:::

## Behavior

```go
expiresOn, err := time.Parse(time.RFC3339, expiresOnStr)
// ...

credential := staticTokenCredential{
  token:     token,
  expiresOn: expiresOn,
}

client, err := azlogs.NewClient(endpoint, credential, &azlogs.ClientOptions{
  ClientOptions: azcore.ClientOptions{
    Transport: &http.Client{
      Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
          InsecureSkipVerify: allowInsecure,
        },
      },
    },
  },
})
// ...

message, err := json.Marshal([]map[string]string{logRecord.Fields})
// ...

_, err = client.Upload(ctx, ruleID, streamName, message, nil)
// ...
```
