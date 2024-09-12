#!/bin/sh

set -e

FLOWG_CMD="docker run --rm -v ./data:/data -p 5080:5080/tcp -p 5514:5514/udp linksociety/flowg:latest"

sudo rm -rf logs.txt data/logs data/auth

${FLOWG_CMD} admin role create --name admin \
  write_streams \
  write_transformers \
  write_pipelines \
  write_acls \
  write_alerts \
  send_logs

${FLOWG_CMD} admin user create --name root --password root admin

export FLOWG_TOKEN=$(
  ${FLOWG_CMD} admin token create --user root
)

${FLOWG_CMD} serve --verbose > logs.txt &
pid=$!
trap "kill $pid" EXIT

echo -n "Waiting for Flowg to start..."
for _ in $(seq 1 10)
do
  echo -n "."
  sleep 0.25
  nc -z localhost 5080 2>/dev/null && echo " ok" && break
done

nc -z localhost 5080 2>/dev/null || (echo " timeout" && exit 1)
