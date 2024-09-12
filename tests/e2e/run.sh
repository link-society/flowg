#!/bin/sh

set -e

. ../start_flowg.sh

hurl --variable token=${FLOWG_TOKEN} --test specs/
