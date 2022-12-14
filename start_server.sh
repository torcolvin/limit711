#!/bin/bash
set -eux

DOCKER_IMAGE=${1:-couchbase/server:enterprise-7.1.0}

# kill couchbase if it exists
docker kill couchbase || true

echo "Starting couchbase server"

COUCHBASE_DATA_DIR=${PWD}/cbs
sudo rm -rf ${COUCHBASE_DATA_DIR}
tar xf cbs-data.tar.bz2
mkdir -p ${COUCHBASE_DATA_DIR}

docker run --rm -d --name couchbase -p 8091-8096:8091-8096 -p 11207:11207 -p 11210:11210 -p 11211:11211 -p 18091-18094:18091-18094 --mount "type=bind,src=${COUCHBASE_DATA_DIR},target=/opt/couchbase/var" $DOCKER_IMAGE

# Test to see if Couchbase Server is up
# Each retry min wait 5s, max 10s. Retry 20 times with exponential backoff (delay 0), fail at 120s
curl --retry-all-errors --connect-timeout 5 --max-time 10 --retry 20 --retry-delay 0 --retry-max-time 120 'http://127.0.0.1:8091'

# Set up CBS
curl -u Administrator:password http://127.0.0.1:8091/nodes/self/controller/settings -d 'path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&'
curl -u Administrator:password http://127.0.0.1:8091/node/controller/setupServices -d 'services=kv%2Cn1ql%2Cindex'
curl -u Administrator:password http://127.0.0.1:8091/pools/default -d 'memoryQuota=3072' -d 'indexMemoryQuota=3072' -d 'ftsMemoryQuota=256'
curl -u Administrator:password http://127.0.0.1:8091/settings/web -d 'password=password&username=Administrator&port=SAME'
curl -u Administrator:password http://localhost:8091/settings/indexes -d indexerThreads=4 -d logLevel=verbose -d maxRollbackPoints=10 \
    -d storageMode=plasma -d memorySnapshotInterval=150 -d stableSnapshotInterval=40000

echo "Finished starting couchbase server"
