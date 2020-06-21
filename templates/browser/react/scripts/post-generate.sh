#!/usr/bin/env sh

set -e
set -u
set -x
set -o pipefail

chmod +x ./scripts/*.sh

yarn
