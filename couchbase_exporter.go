package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/lelvisl/gocbmgr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", ":9131", "Address to listen on for web interface and telemetry.")
	metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	nodeName      = flag.String("node.name", "", "Hostname to filter node metrics.")
	nodeURL       = flag.String("node.url", "http://localhost:8091", "DB Url")
	nodeAuth      = flag.String("node.auth", "", "Couchbase auth - login:password")
)

func main() {

	flag.Parse()
	prometheus.Register(ReplicaNumber)
	prometheus.Register(Stats)
	prometheus.Register(Quota)

	var login, password string
	if len(*nodeAuth) > 0 {
		login = strings.Split(*nodeAuth, ":")[0]
		password = strings.Split(*nodeAuth, ":")[1]
	}
	couchCluster := cbmgr.New([]string{*nodeURL}, login, password)
	// inf,err:=couchCluster.ClusterInfo()
	// if err !=nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%+v\n",inf)

	getStats(couchCluster)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
