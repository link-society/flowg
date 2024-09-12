#!/bin/sh

set -e

. ../flowg.sh

python generate-logs.py --token $FLOWG_ADMIN_TOKEN
