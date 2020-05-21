#!/usr/bin/env bash

set -eux

psql -U postgres -f ./sql/init.sql
PGPASSWORD=todo psql -U todo -f ./sql/schema.sql
