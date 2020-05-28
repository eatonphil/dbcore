#!/usr/bin/env bash

(cd .. && go mod init {{ api.extra.repo }})
goimports -w ./
