#!/bin/sh

set -e

. ../flowg.sh

sudo cp -r config/* data/config/

python -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt

echo "--------------------------------------------------------------------------------"

python generate-logs.py --token $FLOWG_ADMIN_TOKEN
