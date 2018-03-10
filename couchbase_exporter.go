package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lelvisl/couchbase_exporter/version"
	"github.com/lelvisl/gocbmgr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configPath = flag.String("config", "./config.yml", "Config")
	v          = flag.Bool("version", false, "show version")
)

func main() {
	flag.Parse()

	prometheus.Register(replicaNumber)
	prometheus.Register(stats)
	prometheus.Register(quota)
	prometheus.Register(clusterStats)
	prometheus.Register(clusterQuota)

	if *v {
		fmt.Println(version.Show())
		os.Exit(0)
	}
	c, err := configure(*configPath)
	if err != nil {
		log.Printf("Configure err: %s\n", err.Error())
		os.Exit(2)
	}
	var couchCluster *cbmgr.Couchbase
	if len(c.Node.Auth.User) > 0 && len(c.Node.Auth.Password) > 0 {
		couchCluster = cbmgr.New(c.Node.Auth.User, c.Node.Auth.Password)
	} else {
		flag.PrintDefaults()
		os.Exit(254)
	}
	couchCluster.SetEndpoints(c.Node.URLs)
	_, err = couchCluster.ClusterInfo()
	if err != nil {
		log.Printf("can't connect to cluster err: %s\n", err.Error())
		os.Exit(2)
	}

	go func() {
		for {
			getBucketStats(couchCluster)
			getClusterStats(couchCluster)
			time.Sleep(c.Node.Refresh)
		}
	}()

	http.Handle(c.Web.URI, promhttp.Handler())
	server := &http.Server{
		Addr: c.Web.Adress,
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go gracefulShutdown(signals, server)
	log.Println(server.ListenAndServe())
}

func gracefulShutdown(killSignal <-chan os.Signal, s *http.Server) {
	log.Println("graceful shutdown", <-killSignal)
	s.Shutdown(context.Background())
}
