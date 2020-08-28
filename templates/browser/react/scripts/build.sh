#!/usr/bin/env bash

set -eux

# Easier than depending on a yaml parser? TODO: deal with TLS in server
API_PREFIX="http://$(cat ../dbcore.yml | grep address | awk '{ print $2 }' | xargs)"
yarn esbuild src/main.tsx --bundle --define:window.DBCORE_API_PREFIX=\"$API_PREFIX\" '--define:process.env.NODE_ENV="production"' --minify --outfile=build/bundle.js
