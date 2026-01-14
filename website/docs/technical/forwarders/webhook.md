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
  }

  HttpHeader {
    str name "Header name"
    str value "Header value"
  }

  Configuration ||--o{ HttpHeader : has
```

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
