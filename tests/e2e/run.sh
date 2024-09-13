#!/bin/sh

set -e

. ../flowg.sh

hurl \
  --variable admin_token=${FLOWG_ADMIN_TOKEN} \
  --variable guest_token=${FLOWG_GUEST_TOKEN} \
  --report-html reports/html \
  --report-junit reports/junit.xml \
  --test integration/

hurl \
  --variable admin_token=${FLOWG_ADMIN_TOKEN} \
  --variable guest_token=${FLOWG_GUEST_TOKEN} \
  --repeat 1000 \
  --test benchmark/
