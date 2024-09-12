#!/bin/sh

set -e

. ../start_flowg.sh

python generate-logs.py --token $FLOWG_TOKEN
