#!/bin/bash

function flowg_start() {
  echo -n "Starting FlowG..."

  DOCKER_CONTAINER_ID=$(
    docker run -d --rm \
      -v ./data:/data \
      -p 5080:5080/tcp \
      -p 9113:9113/tcp \
      -p 5514:5514/udp \
      linksociety/flowg:latest
  )
  trap "docker kill ${DOCKER_CONTAINER_ID} >/dev/null && flowg_cleanup_data" EXIT

  for _ in $(seq 1 10)
  do
    echo -n "."
    sleep 0.25
    nc -z localhost 9113 2>/dev/null && echo " ok" && break
  done

  nc -z localhost 9113 2>/dev/null || (echo " timeout" && exit 1)
}

function flowg_cleanup_data() {
  echo -n "Kill dangling processes..."
  killall firefox-bin geckodriver 2>/dev/null || true
  echo " ok"

  echo -n "Cleaning up data..."
  sudo rm -rf ./data/
  echo " ok"
}

function flowg_acl_admin() {
  resp=$(
    curl -sf \
      -X POST \
      -H "Content-Type: application/json" \
      -d '{"username":"root","password":"root"}' \
      http://localhost:5080/api/v1/auth/login
  )
  if [[ $? -ne 0 ]]; then exit 11; fi
  jwt=$(echo $resp | jq -r .token)

  resp=$(
    curl -sf \
      -X POST \
      -H "Authorization: Bearer jwt:${jwt}" \
      -H "Content-Type: application/json" \
      http://localhost:5080/api/v1/token
  )
  if [[ $? -ne 0 ]]; then exit 12; fi
  pat=$(echo $resp | jq -r .token)

  echo $pat
}

function flowg_acl_guest() {
  resp=$(
    curl -sf \
      -X PUT \
      -H "Authorization: Bearer pat:${FLOWG_ADMIN_TOKEN}" \
      -H "Content-Type: application/json" \
      -d '{"password":"guest","roles":[]}' \
      http://localhost:5080/api/v1/users/guest
  )
  if [[ $? -ne 0 ]]; then exit 21; fi

  resp=$(
    curl -sf \
      -X POST \
      -H "Content-Type: application/json" \
      -d '{"username":"guest","password":"guest"}' \
      http://localhost:5080/api/v1/auth/login
  )
  if [[ $? -ne 0 ]]; then exit 22; fi
  jwt=$(echo $resp | jq -r .token)

  resp=$(
    curl -sf \
      -X POST \
      -H "Authorization: Bearer jwt:${jwt}" \
      -H "Content-Type: application/json" \
      http://localhost:5080/api/v1/token
  )
  if [[ $? -ne 0 ]]; then exit 23; fi
  pat=$(echo $resp | jq -r .token)

  echo $pat
}

function _error() {
  case "$1" in
    11) echo "Failed to authenticate as admin" ;;
    12) echo "Failed to generate admin token" ;;
    21) echo "Failed to create guest user" ;;
    22) echo "Failed to authenticate as guest" ;;
    23) echo "Failed to generate guest token" ;;
    *) echo "Unknown error: $1" ;;
  esac
}

echo "--------------------------------------------------------------------------------"

flowg_cleanup_data

flowg_start

echo -n "Setup ACLs..."
FLOWG_ADMIN_TOKEN=$(flowg_acl_admin)
code=$?; if [[ $code -ne 0 ]]; then echo " error: $(_error $code)"; exit 1; fi

FLOWG_GUEST_TOKEN=$(flowg_acl_guest)
code=$?; if [[ $code -ne 0 ]]; then echo " error: $(_error $code)"; exit 1; fi

export FLOWG_ADMIN_TOKEN
export FLOWG_GUEST_TOKEN

echo " ok"

echo "--------------------------------------------------------------------------------"
