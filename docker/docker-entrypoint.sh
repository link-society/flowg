#!/bin/sh

set -e

mkdir -p /data
chmod 0700 /data
chown -R flowg:flowg /data

exec su-exec flowg /usr/local/bin/flowg-server $@
