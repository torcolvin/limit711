#!/bin/bash
set -eux -o pipefail

/bin/bash ./start_server.sh couchbase/server:enterprise-7.1.1
go run .
