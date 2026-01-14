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
  }
```

*Notes:*

1. The token is **NOT** encrypted in the database.

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
  "source": "flowg",
  "host": "<log.host>",
  "time": "<timestamp>",
}
```
