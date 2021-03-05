package main

import (
	"log"

	"github.com/lelvisl/gocbmgr"
)

func getBucketStats(couch *cbmgr.Couchbase) {
	buckets, err := couch.GetBuckets()
	if err != nil {
		log.Println(err)
	}
	for _, bucket := range buckets {
		stat, err := couch.GetBucketStatus(bucket.BucketName)
		if err != nil {
			log.Println(err)
			continue
		}
		for k, v := range stat.Quota {
			quota.WithLabelValues(bucket.BucketName, k).Set(float64(v))
		}
		replicaNumber.WithLabelValues(bucket.BucketName).Set(float64(stat.ReplicaNumber))
		//example - пример вывода stats
		//TODO: подумать, как это завернуть в метрики прома. Пока думаю, что надо прать первый семпл, и отдавать его как метрику за минуту. (учитывая что семплы мы за минуту и выбираем по дефолту)
		monStats, err := couch.GetBucketStats(bucket.BucketName)
		if err != nil {
			log.Println(err)
			continue
		}
		for itemName, v := range monStats {
			for _, item := range v.Value {
				// TODO: нужно проверять числа вида 2.28170137e+08 и приводить их к сайзингу
				val := item[len(item)-1]
				stats.WithLabelValues(bucket.BucketName, itemName).Set(val)
			}
		}
	}

}

func getClusterStats(couch *cbmgr.Couchbase) {
	clstInf, err := couch.ClusterInfo()
	if err != nil {
		log.Printf("getClusterStats error:%s", err.Error())
	}
	if clstInf.Balanced {
		clusterStats.WithLabelValues("balanced").Set(0)
	} else {
		clusterStats.WithLabelValues("balanced").Set(1)
	}
	//RebalaceStatus ????
	clusterQuota.WithLabelValues("DataMemoryQuotaMB").Set(float64(clstInf.DataMemoryQuotaMB))
	clusterQuota.WithLabelValues("IndexMemoryQuotaMB").Set(float64(clstInf.IndexMemoryQuotaMB))
	clusterQuota.WithLabelValues("SearchMemoryQuotaMB").Set(float64(clstInf.SearchMemoryQuotaMB))
}
