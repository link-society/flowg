#!/bin/sh

rm -rf logs.txt data/logs

../../bin/flowg \
    -log-dir ./data/logs \
    -config-dir ./data/config \
    -bind 127.0.0.1:5080 \
    -verbose \
  > logs.txt &
pid=$!
trap "kill $pid" EXIT

sleep 0.1

python generate-logs.py
