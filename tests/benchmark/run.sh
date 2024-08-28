#!/bin/sh

rm -rf logs.txt data/logs data/auth

../../bin/flowg admin role create \
  --auth-dir ./data/auth \
  --name admin \
  write_streams \
  write_transformers \
  write_pipelines \
  write_acls \
  send_logs

../../bin/flowg admin user create \
  --auth-dir ./data/auth \
  --name root \
  --password root \
  admin

token=$(
  ../../bin/flowg admin token create \
    --auth-dir ./data/auth \
    --user root
)

../../bin/flowg serve \
    --log-dir ./data/logs \
    --config-dir ./data/config \
    --bind 127.0.0.1:5080 \
    --verbose \
  > logs.txt &
pid=$!
trap "kill $pid" EXIT

sleep 0.1

python generate-logs.py --token $token
