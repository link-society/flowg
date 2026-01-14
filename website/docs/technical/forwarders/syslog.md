---
sidebar_position: 2
---

# Syslog

This forwarder is used to send a log record to a remote Syslog Server.

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str network "Either TCP or UDP"
    str address "The `host:port` address of the Syslog Server"
    str tag[1] "The tag to use"
    str severity[1] "The severity level to use"
    str facility[1] "The facility to log to"
    str message[1] "The content of the message to log"
  }
```

*Notes:*

1. Those fields are "dynamic", consult
[this page](/docs/technical/dynamic-fields) for more information.

## Behavior

```go
priority := severity | facility
writer, err := syslog.Dial(network, address, priority, tag)
// ...
err = writer.Write(message)
// ...
```
