---
sidebar_position: 1
---

# Webhook

This forwarder is used to send a log record to an HTTP endpoint using a POST
request.

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str url "The target URL"
    HttpHeader[] headers "Additional HTTP headers to send"
    str body[1] "The body of a log message"
  }

  HttpHeader {
    str name "Header name"
    str value "Header value"
  }

  Configuration ||--o{ HttpHeader : has
```

*Notes: *

1. This field is "dynamic", consult
   [this page](/docs/technical/dynamic-fields) for more information.

## Behavior

```
POST http://example.com
{
  "timestamp": "...",
  "fields": {
    "...": "..."
  }
}
```
