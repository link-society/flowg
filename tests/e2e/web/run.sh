#!/bin/bash

. ../../flowg.sh

set -e

if [ ! -d "venv" ]; then
  python3 -m venv venv
fi

. venv/bin/activate

sudo apt install libasound2t64  # required for Firefox webdriver
pip install -r requirements.txt

rm -rf reports/
export ROBOT_OPTIONS="--outputdir reports/"
robot spec/
