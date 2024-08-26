#!/bin/sh

case "$1" in
  serve)
    shift
    exec /usr/local/bin/flowg serve $@
    ;;

  admin)
    shift
    exec /usr/local/bin/flowg admin $@
    ;;

  *)
    exec /usr/local/bin/flowg --help
    ;;
esac
