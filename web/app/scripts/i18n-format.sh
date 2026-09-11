#!/bin/sh

set -ex

msgcat --sort-output --no-wrap src/locales/dummy.po -o src/locales/dummy.po
msgcat --sort-output --no-wrap src/locales/en.po -o src/locales/en.po
