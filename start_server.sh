#!/bin/bash

set -eux

DOCKER_IMAGE=couchbase/server:enterprise-7.1.1

# kill couchbase if it exists
docker kill couchbase || true
docker volume rm couchbase || true

docker run --rm -d --name couchbase -p 8091-8096:8091-8096 -p 11207:11207 -p 11210:11210 -p 11211:11211 -p 18091-18094:18091-18094 --volume 'type=volume,src=couchbase_data,/var/couchbase/var' $DOCKER_IMAGE

# Test to see if Couchbase Server is up
# Each retry min wait 5s, max 10s. Retry 20 times with exponential backoff (delay 0), fail at 120s
curl --retry-all-errors --connect-timeout 5 --max-time 10 --retry 20 --retry-delay 0 --retry-max-time 120 'http://127.0.0.1:8091'

# Set up CBS
curl -u Administrator:password -v -X POST http://127.0.0.1:8091/nodes/self/controller/settings -d 'path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&'
curl -u Administrator:password -v -X POST http://127.0.0.1:8091/node/controller/setupServices -d 'services=kv%2Cn1ql%2Cindex'
curl -u Administrator:password -v -X POST http://127.0.0.1:8091/pools/default -d 'memoryQuota=3072' -d 'indexMemoryQuota=3072' -d 'ftsMemoryQuota=256'
curl -u Administrator:password -v -X POST http://127.0.0.1:8091/settings/web -d 'password=password&username=Administrator&port=SAME'
curl -u Administrator:password -v -X POST http://localhost:8091/settings/indexes -d indexerThreads=4 -d logLevel=verbose -d maxRollbackPoints=10 \
    -d storageMode=plasma -d memorySnapshotInterval=150 -d stableSnapshotInterval=40000

echo ""

docker cp archives couchbase:/archives

docker exec couchbase cbbackupmgr restore --archive /archives --repo example --cluster couchbase://127.0.0.1 --username Administrator --password password --auto-create-buckets

