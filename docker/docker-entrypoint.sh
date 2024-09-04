#!/bin/sh

set -e

mkdir -p /data
chmod 0700 /data
chown -R flowg:flowg /data

case "$1" in
  serve)
    shift
    exec su-exec flowg /usr/local/bin/flowg serve $@
    ;;

  admin)
    shift
    exec su-exec flowg /usr/local/bin/flowg admin $@
    ;;

  *)
    exec su-exec flowg /usr/local/bin/flowg --help
    ;;
esac
