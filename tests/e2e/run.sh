#!/bin/sh

set -e

. ../flowg.sh

hurl \
  --variable admin_token=${FLOWG_ADMIN_TOKEN} \
  --variable guest_token=${FLOWG_GUEST_TOKEN} \
  --test specs/
