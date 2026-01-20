---
sidebar_position: 4
---

# Splunk

This forwarder is used to send a log record to [Splunk](http://splunk.com/).

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str endpoint "URL of the Splunk HTTP Event Collector endpoint"
    str token[1] "Token used for authentication"
    str source[2] "Source of Log (defaults to 'flowg')"
    str host[2] "The host to be sent to Splunk"
  }
```

*Notes:*

1. The token is **NOT** encrypted in the database.
2. These fields are "dynamic", consult
   [this page](/docs/technical/dynamic-fields) for more information.

## Behavior

```
POST <ENDPOINT>
Authorization: Splunk <TOKEN>
Content-Type: application/json

{
  "event": {
    "...": "...",
  },
  "sourcetype": "json",
  "source": "<log.source>",
  "host": "<log.host>",
  "time": "<timestamp>",
}
```
