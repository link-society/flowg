---
sidebar_position: 5
---

# AMQP

This forwarder is used to send a log record to an AMQP broker, such as
[RabbitMQ](http://rabbitmq.com/).

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str url "URL of AMQP broker"
    str exchange "Name of the exchange to send the log to"
    str routing_key "Routing Key to use when sending the log to the exchange"
  }
```

*Notes:*

1. The token is **NOT** encrypted in the database.

## Behavior

```go
payload, err := json.Marshal(logRecord)

conn, err := amqp.Dial(url)
ch, err := conn.Channel()
err = ch.Publish(
  exchange,
  routingKey,
  /* mandatory= */ false,
  /* immediate= */ false,
  amqp.Publishing{
    ContentType: "application/json",
    Body:        payload,
  },
)
```
