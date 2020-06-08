#!/usr/bin/env bash

set -eux

yarn esbuild src/main.tsx --bundle '--define:process.env.NODE_ENV="production"' --minify --outfile=build/bundle.js
