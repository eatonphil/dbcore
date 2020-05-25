#!/usr/bin/env bash

psql -U postgres -c 'DROP DATABASE todo; DROP USER todo;'
