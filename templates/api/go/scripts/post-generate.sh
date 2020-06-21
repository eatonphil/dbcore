#!/usr/bin/env bash

set -eu

go get golang.org/x/tools/cmd/goimports
if [[ ! -f ../go.mod ]]; then
    (cd .. && go mod init {{ api.extra.repo }})
fi
goimports -w ./
