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
    str source[2] "The integration name associated with your log"
    str tags[2] "Tags associated with your logs"
    str hostname[2] "The name of the originating host of the log"
    str message[2] "The message reserved attribute of your log"
    str service[2] "The name of the application or service generating the log events"
  }
```

*Notes:*

1. The key is **NOT** encrypted in the database.
2. Those fields are "dynamic", consult
   [this page](/docs/technical/dynamic-fields) for more information.
3. For more information about Datadog parameters visit [this page](https://docs.datadoghq.com/api/latest/logs/#send-logs)

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
