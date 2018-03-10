package main

import "github.com/prometheus/client_golang/prometheus"

var (
	replicaNumber = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "replica",
		Help:      "replica",
	},
		[]string{"bucket"},
	)
	quota = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "quota",
		Help:      "quota",
	},
		[]string{"bucket", "type"},
	)

	stats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "stats",
		Help:      "see /pools/default/buckets/Data/statsDirectory",
	},
		[]string{"bucket", "item"},
	)
	clusterStats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "cluster",
		Name:      "stats",
		Help:      "will be soon",
	},
		[]string{"item"},
	)
	clusterQuota = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "cluster",
		Name:      "quota",
		Help:      "will be soon",
	},
		[]string{"type"},
	)
)
