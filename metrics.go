package main

import "github.com/prometheus/client_golang/prometheus"

var (
	ReplicaNumber = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "replica",
		Help:      "replica",
	},
		[]string{"bucket"},
	)
	Quota = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "quota",
		Help:      "quota",
	},
		[]string{"bucket", "type"},
	)

	Stats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "stats",
		Help:      "see /pools/default/buckets/Data/statsDirectory",
	},
		[]string{"bucket", "item"},
	)
	ClusterStats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "cluster",
		Name:      "stats",
		Help:      "will be soon",
	},
		[]string{"item"},
	)
	ClusterQuota = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "cluster",
		Name:      "quota",
		Help:      "will be soon",
	},
		[]string{"type"},
	)
)
