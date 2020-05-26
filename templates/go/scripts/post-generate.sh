#!/usr/bin/env bash

(cd .. && go mod init {{ api.repo }})
goimports -w ./
