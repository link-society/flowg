#!/bin/sh
set -eu

PORT="${PORT:-8080}"

socat TCP-LISTEN:"$PORT",reuseaddr,fork SYSTEM:"/dns-management-router.sh"
