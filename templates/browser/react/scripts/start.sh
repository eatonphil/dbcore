#!/usr/bin/env bash

set -eux

rm -rf build
mkdir build
yarn esbuild src/main.tsx --bundle '--define:process.env.NODE_ENV="development"' --outfile=build/bundle.js
cp -r static/* build/
yarn es-dev-server --port 9091 --root-dir build --app-index index.html
