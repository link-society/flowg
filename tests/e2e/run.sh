#!/bin/sh

set -e

. ../flowg.sh

hurl \
  --variable admin_token=${FLOWG_ADMIN_TOKEN} \
  --variable guest_token=${FLOWG_GUEST_TOKEN} \
  --variable timewindow_begin=$(date -d "5 minutes ago" -u +"%Y-%m-%dT%H:%M:%SZ") \
  --variable timewindow_end=$(date -d "+5 minutes" -u +"%Y-%m-%dT%H:%M:%SZ") \
  --error-format long \
  --report-html reports/html \
  --report-junit reports/junit.xml \
  --jobs 1 \
  --test spec/
