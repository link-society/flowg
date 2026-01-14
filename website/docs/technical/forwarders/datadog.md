---
sidebar_position: 3
---

# Datadog

This forwarder is used to send a log record to [Datadog](http://datadoghq.com/).

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str url "URL of the Datadog Intake endpoint"
    str apiKey[1] "API key used for authentication"
  }
```

*Notes:*

1. The key is **NOT** encrypted in the database.

## Behavior

```
POST <URL>
Accept: application/json
Content-Type: application/json
DD-API-KEY: <APIKEY>

{
  "ddtags": "<log.ddtags>",
  "ddsource": "<log.ddsource>",
  "hostname": "<log.hostname>",
  "service": "<log.service>",
  "message": "<log.message>",
}
```
