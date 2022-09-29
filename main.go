package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

const (
	withLimit = "LIMIT 5" // turn to empty string to avoid error
	//withLimit   = "" // turn to empty string to avoid error
	cbsUserName = "Administrator"
	cbsPassword = "password"
	cbsAddr     = "localhost"
	bucketName  = "sg_int_0"
	//idxName     = "sg_syncDocs_1"
)

func main() {
	log.Print("Start test program")
	cluster := getCluster()
	//buildIndexes(cluster)
	query(cluster)

	log.Print("SUCCESS")
}

type QueryIdRow struct {
	Id string
}

func query(cluster *gocb.Cluster) {
	keyspace := fmt.Sprintf("`%s`.`_default`.`_default`", bucketName)
	statement := fmt.Sprintf("SELECT META(sgQueryKeyspaceAlias).id FROM %s AS sgQueryKeyspaceAlias USE INDEX(sg_syncDocs_1) WHERE META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:%%' AND (META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:user:%%' OR META(sgQueryKeyspaceAlias).id LIKE '\\\\_sync:role:%%') AND META(sgQueryKeyspaceAlias).id >= $startkey ORDER BY META(sgQueryKeyspaceAlias).id %s", keyspace, withLimit)

	fmt.Printf("%q\n", statement)
	results, err := cluster.Query(statement, &gocb.QueryOptions{
		ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
		NamedParameters: map[string]interface{}{
			"startkey": "",
		},
	})
	if err != nil {
		log.Fatalf("Query returned error: %s", err)
	}
	var users []string
	for {
		found := results.Next()
		if !found {
			break
		}
		var queryRow QueryIdRow
		err := results.Row(&queryRow)
		if err != nil {
			log.Fatalf("Could not read row: %s", err)
		}
		users = append(users, queryRow.Id)
	}
	fmt.Printf("users=%s\n", users)

}

func getCluster() *gocb.Cluster {
	DefaultGocbV2OperationTimeout := 10 * time.Second

	clusterOptions := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: cbsUserName,
			Password: cbsPassword,
		},
		SecurityConfig: gocb.SecurityConfig{
			TLSSkipVerify: false,
		},
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout:    DefaultGocbV2OperationTimeout,
			KVTimeout:         DefaultGocbV2OperationTimeout,
			ManagementTimeout: DefaultGocbV2OperationTimeout,
			QueryTimeout:      90 * time.Second,
			ViewTimeout:       90 * time.Second,
		},
	}
	connStr := fmt.Sprintf("couchbase://%s?idle_http_connection_timeout=90000&kv_pool_size=2&max_idle_http_connections=64000&max_perhost_idle_http_connections=256", cbsAddr)
	cluster, err := gocb.Connect(connStr, clusterOptions)
	if err != nil {
		log.Fatalf("Error connecticting to cluster %+v", err)
	}
	err = cluster.WaitUntilReady(15*time.Second, nil)
	if err != nil {
		log.Fatalf("Can't connect to cluster %+v", err)
	}

	err = cluster.WaitUntilReady(90*time.Second,
		&gocb.WaitUntilReadyOptions{ServiceTypes: []gocb.ServiceType{gocb.ServiceTypeQuery}},
	)
	if err != nil {
		log.Fatalf("Query service not online")
	}
	return cluster
}

/*
func buildIndexes(cluster *gocb.Cluster) {
	existingIdxs, err := cluster.QueryIndexes().GetAllIndexes(bucketName, nil)
	if err != nil {
		log.Fatalf("Could not GetAllIndexes: %s", err)
	}
	if len(existingIdxs) != 0 {
		if len(existingIdxs) > 1 {
			log.Fatalf("Got too many indexes %+v", existingIdxs)
		}
		if existingIdxs[0].Name != idxName {
			log.Fatalf("Expected %+v but got %+v", idxName, existingIdxs[0].Name)
		}
		return
	}
	syncPrefix := `"\\_sync%"`
	deferBuild := `WITH {  "defer_build":true }`
	statement := fmt.Sprintf("CREATE INDEX `%s` ON `sg_int_0`((meta().`id`)) WHERE ((meta().`id`) like %s) %s", idxName, syncPrefix, deferBuild)
	results, err := cluster.Query(statement, &gocb.QueryOptions{
		ScanConsistency: gocb.QueryScanConsistencyRequestPlus,
	},
	)
	if err != nil {
		log.Fatalf("Could not create index: %s", err)
	}
	// Drain results to return any non-query errors
	for results.Next() {
	}
	closeErr := results.Close()
	if closeErr != nil {
		log.Fatalf("Could not close query: %s", closeErr)

	}
	idx, err := cluster.QueryIndexes().BuildDeferredIndexes(bucketName, nil)
	if err != nil {
		log.Fatalf("Could not build deferred index: %s", err)
	}
	if len(idx) != 1 {
		log.Fatalf("Unexpected numbers of indexes %+v", idx)
	}
	err = cluster.QueryIndexes().WatchIndexes(bucketName, []string{idxName}, 30*time.Second, nil)
	if err != nil {
		log.Fatalf("Index did not come online")
	}

}
*/
