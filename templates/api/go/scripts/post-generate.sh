#!/usr/bin/env bash

go get golang.org/x/tools/cmd/goimports
(cd .. && go mod init {{ api.extra.repo }})
goimports -w ./
