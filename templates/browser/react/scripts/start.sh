#!/usr/bin/env bash

set -eux

rm -rf build
mkdir build
# Easier than depending on a yaml parser? TODO: deal with TLS in server
API_PREFIX="http://$(cat ../api.yml | grep address | cut -d ' ' -f 2 | xargs)"
yarn esbuild src/main.tsx --bundle --define:window.DBCORE_API_PREFIX=\"$API_PREFIX\" '--define:process.env.NODE_ENV="development"' --outfile=build/bundle.js
cp -r static/* build/
yarn es-dev-server --port 9091 --root-dir build --app-index index.html
