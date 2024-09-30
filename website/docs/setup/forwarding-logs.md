---
sidebar_position: 2
---

# Forwarding logs

*FlowG* provides an UDP endpoint capable of receiving Syslog events.

The event will be sent to all pipelines, it is up to the user to filter out the
events from the `SYSLOG` source node.

### In Kubernetes

The Helm chart deploys [Fluentd](https://www.fluentd.org) alongside *FlowG* to
automatically forward the logs of every pod. No further configuration is
required.

### Using Docker

Configure the default log-driver in `/etc/docker/daemon.json`:

```json
{
  "log-driver": "syslog",
  "log-opts": {
    "syslog-address": "udp://127.0.0.1:5514"
  }
}
```

> **NB:** Changing the default logging driver or logging driver options in the
> daemon configuration only affects containers that are created after the
> configuration is changed. Existing containers retain the logging driver
> options that were used when they were created. To update the logging driver
> for a container, the container has to be re-created with the desired options.

### Using syslog-ng

In `/etc/syslog-ng/syslog-ng.conf`:

```
destination d_flowg {
  udp("127.0.0.1" port(5514))
}
```

### Using rsyslog

In `/etc/rsyslog.conf`:

```
*.* @127.0.0.1:5514
```

### Using Logstash (with the Syslog output plugin)

Install the Syslog output plugin by following
[those instructions](https://www.elastic.co/guide/en/logstash/current/plugins-outputs-syslog.html).

Then, in `/etc/logstash/conf.d/flowg.conf`:

```
output {
  syslog {
    host => "127.0.0.1"
    port => 5514
    protocol => "udp"
  }
}
```
