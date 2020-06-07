#!/usr/bin/env bash

set -eux

rm -rf build
mkdir build
yarn esbuild src/main.tsx --bundle '--define:process.env.NODE_ENV="development"' --outfile=build/bundle.js
cp -r static/* build/
python3 -m http.server --directory build 9091
