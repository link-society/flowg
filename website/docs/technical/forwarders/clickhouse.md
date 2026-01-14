---
sidebar_position: 8
---

# Clickhouse

This forwarder is used to persist a log record in an
[Clickhouse](https://clickhouse.com) table.

> **NB:** If the target table does not exist, it will be created.

## Data Model

```mermaid
erDiagram
  direction LR

  Configuration {
    str address "`host:port` of the Clickhouse server to connect to"
    str username "Username to use to authenticate to the Clickhouse server"
    str password[1] "Password to use to authenticate to the Clickhouse server"
    bool tls "Define if TLS handshake is required"
    str db "Name of the database to store the log in"
    str table "Name of the table to store the log in"
  }
```

*Notes:*

1. The password is **NOT** encrypted in the database.

## Behavior

```sql
CREATE TABLE IF NOT EXISTS `table` (
	id         UUID                 NOT NULL PRIMARY KEY,
	timestamp  DateTime64(3, 'UTC') NOT NULL,
	fields     Map(String, String)  NOT NULL,
) ENGINE = MergeTree;

INSERT INTO `table`
VALUES (`uuidv4`, `logRecord.timestamp`, `logRecord.fields`);
```
