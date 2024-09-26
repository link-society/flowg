---
sidebar_position: 4
---

# Alerting

Alerts are HTTP(S) Webhooks that a pipeline can call to send the log to a third
party system like [Zapier](https://zapier.com).

When the alert is triggered, *FlowG* will send a `POST` request to the webhook
URL, with the webhook's HTTP headers. The body of the request will be the log
record itself:

```
POST http://example.com/
Authorization: Bearer s3cr3!

{"timestamp": "2024-01-01T20:00:00.000Z", "field": {"message": "hello"}}
```
