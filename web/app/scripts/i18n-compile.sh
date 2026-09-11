#!/bin/sh

set -ex

mkdir -p src/locales/generated/dummy
i18next-conv -l dummy -s src/locales/dummy.po -t src/locales/generated/dummy/translation.json

mkdir -p src/locales/generated/en
i18next-conv -l en -s src/locales/en.po -t src/locales/generated/en/translation.json
