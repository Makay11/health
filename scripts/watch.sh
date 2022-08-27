#!/usr/bin/env bash

# workaround to get sudo permissions before running gow
# gow does not work properly when the exec script asks for the sudo password
sudo printf ''

go run github.com/mitranim/gow -c run -exec ./tools/setcap.sh $@
