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
	Version        = flag.Bool("version", false, "show version")
	ConfigFileName = flag.String("config", "config.toml", "config file in toml format")
)

var (
	// Config
	Configuration Config
)

func main() {
	var (
		couchCluster *cbmgr.Couchbase
	)
	if ok, err := configure("init"); !ok {
		log.Println(err)
		os.Exit(-5)
	}
	prometheus.Register(ReplicaNumber)
	prometheus.Register(Stats)
	prometheus.Register(Quota)
	prometheus.Register(ClusterStats)
	prometheus.Register(ClusterQuota)

	if *Version {
		fmt.Println(version.Show())
		os.Exit(0)
	}

	log.Println(Configuration)
	//if len(*nodeAuth) > 0 {
	//login = strings.Split(*nodeAuth, ":")[0]
	//password = strings.Split(*nodeAuth, ":")[1]
	if len(Configuration.Core.Username) != 0 || len(Configuration.Core.Password) != 0 {
		couchCluster = cbmgr.New(
			Configuration.Core.Username,
			Configuration.Core.Password,
		)
		couchCluster.SetEndpoints(Configuration.Core.NodeURL)
	} else {
		flag.PrintDefaults()
		os.Exit(254)
	}
	go func() {
		for {
			getBucketStats(couchCluster)
			getClusterStats(couchCluster)
			//тут надо добавить duration снаружи, что бы указать, как часто опрашивать кластер
			time.Sleep(Configuration.Core.RefreshInterval.Duration)
		}
	}()

	http.Handle(Configuration.Core.MetricUri, promhttp.Handler())
	server := &http.Server{
		Addr: Configuration.Core.getAddress(),
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
