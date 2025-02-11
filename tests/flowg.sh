#!/bin/sh

FLOWG_CMD_FG="docker run    --rm -v ./data:/data -p 5080:5080/tcp -p 5514:5514/udp linksociety/flowg:latest"
FLOWG_CMD_BG="docker run -d --rm -v ./data:/data -p 5080:5080/tcp -p 5514:5514/udp linksociety/flowg:latest"

FLOWG_CMD_FAILED="docker run -d --rm -e FLOWG_SYSLOG_TLS_ENABLED=true -v ./data:/data -p 5080:5080/tcp -p 5514:5514/udp linksociety/flowg:latest"

flowg_cleanup_data() {
  echo -n "Cleaning up data..."
  sudo rm -rf ./data/
  echo " ok"
}

flowg_acl_admin() {
  ${FLOWG_CMD_FG} admin role create --name admin \
      write_streams \
      write_transformers \
      write_pipelines \
      write_acls \
      write_alerts \
      send_logs \
    > /dev/null

  ${FLOWG_CMD_FG} admin user create --name root --password root admin \
    > /dev/null

  ${FLOWG_CMD_FG} admin token create --user root
}

flowg_acl_guest() {
  ${FLOWG_CMD_FG} admin user create --name guest --password guest \
    > /dev/null

  ${FLOWG_CMD_FG} admin token create --user guest
}

flowg_start() {
  echo -n "Check FlowG Docker image exists..."

  DOCKER_CONTAINER_ID=$(${FLOWG_CMD_FAILED} serve)
  sleep 3

  echo " ok"

  echo -n "Starting FlowG..."

  DOCKER_CONTAINER_ID=$(${FLOWG_CMD_BG} serve)
  trap "docker kill ${DOCKER_CONTAINER_ID} >/dev/null && flowg_cleanup_data" EXIT

  for _ in $(seq 1 10)
  do
    echo -n "."
    sleep 0.25
    nc -z localhost 5080 2>/dev/null && echo " ok" && break
  done

  nc -z localhost 5080 2>/dev/null || (echo " timeout" && exit 1)
}


echo "--------------------------------------------------------------------------------"

flowg_cleanup_data

echo -n "Setup ACLs..."
export FLOWG_ADMIN_TOKEN=$(flowg_acl_admin)
export FLOWG_GUEST_TOKEN=$(flowg_acl_guest)
echo " ok"

flowg_start

echo "--------------------------------------------------------------------------------"
