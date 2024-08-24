#!/bin/sh

exec /usr/local/bin/flowg \
  -bind ":5080" \
  -log-dir /data/logs \
  -config-dir /data/config
