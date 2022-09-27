package main

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

const (
	cbsUserName = "Administrator"
	cbsPassword = "password"
	cbsAddr     = "localhost"
)

func main() {
	log.Print("Start test program")
	getCluster()
	log.Print("SUCCESS")
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
	return cluster
}
