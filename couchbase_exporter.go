package main

import (
	"flag"
	"fmt"
	"github.com/lelvisl/couchbase_exporter/version"
	"github.com/lelvisl/gocbmgr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	listenAddress = flag.String("web.listen-address", ":9131", "Address to listen on for web interface and telemetry.")
	metricUri     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	nodeName      = flag.String("node.name", "", "Hostname to filter node metrics.")
	nodeURL       = flag.String("node.url", "http://localhost:8091", "DB Url")
	nodeAuth      = flag.String("node.auth", "", "Couchbase auth - login:password")
	Version       = flag.Bool("version", false, "show version")
)

func main() {
	var login, password string
	flag.Parse()
	prometheus.Register(ReplicaNumber)
	prometheus.Register(Stats)
	prometheus.Register(Quota)
	if *Version {
		fmt.Println(version.Show())
		os.Exit(0)
	}

	if len(*nodeAuth) > 0 {
		login = strings.Split(*nodeAuth, ":")[0]
		password = strings.Split(*nodeAuth, ":")[1]
	} else {
		flag.PrintDefaults()
		os.Exit(254)
	}
	couchCluster := cbmgr.New([]string{*nodeURL}, login, password, nil)
	getStats(couchCluster)

	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr: *listenAddress,
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
