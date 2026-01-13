---
sidebar_position: 3
---

# Authentication

You can enable authentication between nodes by using a secret key that each node
requires:

```bash
cookie=$(openssl rand -hex 32)

flowg-server \
  --cluster-node-id flowg-node0 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --auth-dir ./data/node0/auth \
  --log-dir ./data/node0/logs \
  --config-dir ./data/node0/config \
  --cluster-state-dir ./data/node0/state \
  --http-bind 127.0.0.1:5080 \
  --mgmt-bind 127.0.0.1:9113 \
  --syslog-bind 127.0.0.1:5514 &

flowg-server \
  --cluster-node-id flowg-node1 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node0 \
  --cluster-formation-manual-join-endpoint http://localhost:9113 \
  --auth-dir ./data/node1/auth \
  --log-dir ./data/node1/logs \
  --config-dir ./data/node1/config \
  --cluster-state-dir ./data/node1/state \
  --http-bind 127.0.0.1:5081 \
  --mgmt-bind 127.0.0.1:9114 \
  --syslog-bind 127.0.0.1:5515 &

flowg-server \
  --cluster-node-id flowg-node2 \
  --cluster-cookie ${cookie} \
  --cluster-formation-strategy="manual" \
  --cluster-formation-manual-join-node-id flowg-node1 \
  --cluster-formation-manual-join-endpoint http://localhost:9114 \
  --auth-dir ./data/node2/auth \
  --log-dir ./data/node2/logs \
  --config-dir ./data/node2/config \
  --cluster-state-dir ./data/node2/state \
  --http-bind 127.0.0.1:5082 \
  --mgmt-bind 127.0.0.1:9115 \
  --syslog-bind 127.0.0.1:5516 &
```
