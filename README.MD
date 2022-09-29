This issue came from Sync Gateway integration tests, and only occurs using Couchbase Server 7.1.1. It does not occur on 7.1.0 or 7.0.3.

The database is a docker snapshot of that created by the integration tests. We use this tarball to re-run this against docker versions.

This reproduction creates a docker image against `couchbase/enterprise-7.1.1` or `couchbase-7.1.0`.

Running `./run_7.1.1.sh` will produce the malformed output:

```
ERROR got not enough users=[_sync:user: _sync:user:user1 _sync:user:user10]
```

Running `./run_7.1.0.sh` against the same database state (from tar file) will produce the expected output:

```
Got expected users= [_sync:user: _sync:user:user1 _sync:user:user10 _sync:user:user11 _sync:user:user12]
```

If you run against 7.1.1, and open in couchbase UI and run the query against the workbench, it will "fix" it for the go program.

```
SELECT META(sgQueryKeyspaceAlias).id FROM `sg_int_0`.`_default`.`_default` AS sgQueryKeyspaceAlias USE INDEX(sg_syncDocs_1) WHERE META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:%' AND (META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:user:%' OR META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:role:%') AND META(sgQueryKeyspaceAlias).id >= $startkey ORDER BY META(sgQueryKeyspaceAlias).id LIMIT 5"
```

Regardless of the number of times you run query in go, i.e. repeat `go run .`, you will get the same malformed result in go, until either you run it in the workbench, or you run an EXPLAIN first.
