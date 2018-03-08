package main

import "github.com/prometheus/client_golang/prometheus"

var (
	ReplicaNumber = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "ReplicaNumber",
		Help:      "ReplicaNumber",
	},
		[]string{"bucket"},
	)
	Quota = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "Quota",
		Help:      "Quota",
	},
		[]string{"bucket", "type"},
	)

	Stats = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "couchbase",
		Subsystem: "bucket",
		Name:      "Stats",
		Help:      "see /pools/default/buckets/Data/statsDirectory",
	},
		[]string{"bucket", "item"},
	)
)
