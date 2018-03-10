package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/lelvisl/couchbase_exporter/version"
	"github.com/lelvisl/gocbmgr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configPath = flag.String("config", "./config.yml", "Config")
	Version    = flag.Bool("version", false, "show version")
)

func main() {
	var login, password string
	flag.Parse()

	prometheus.Register(ReplicaNumber)
	prometheus.Register(Stats)
	prometheus.Register(Quota)
	prometheus.Register(ClusterStats)
	prometheus.Register(ClusterQuota)

	if *Version {
		fmt.Println(version.Show())
		os.Exit(0)
	}
	c, err := configure(*configPath)
	if err != nil {
		log.Println("Configure err: %s", err.Error())
		os.Exit(2)
	}

	if len(c.node.auth) > 0 {
		login = strings.Split(c.node.auth, ":")[0]
		password = strings.Split(c.node.auth, ":")[1]
	} else {
		flag.PrintDefaults()
		os.Exit(254)
	}
	couchCluster := cbmgr.New(login, password)
	couchCluster.SetEndpoints(c.node.urls)
	go func() {
		for {
			getBucketStats(couchCluster)
			getClusterStats(couchCluster)
			//тут надо добавить duration снаружи, что бы указать, как часто опрашивать кластер
			time.Sleep(5 * time.Second)
		}
	}()

	http.Handle(c.web.metricURI, promhttp.Handler())
	server := &http.Server{
		Addr: c.web.listenAddress,
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
