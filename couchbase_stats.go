package main

import (
	"fmt"

	"github.com/lelvisl/gocbmgr"
	"log"
)

func getStats(couch *cbmgr.Couchbase) {
	buckets, err := couch.GetBuckets()
	if err != nil {
		fmt.Println(err)
	}
	for _, bucket := range buckets {
		stat, err := couch.GetBucketStatus(bucket.BucketName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for k, v := range stat.Quota {
			// TODO думаю нужно в конфиге или в аргументах передавать сайзинг b/kb/mb/gb
			Quota.WithLabelValues(bucket.BucketName, k).Set(float64(v) / 1024 / 1024)
		}
		ReplicaNumber.WithLabelValues(bucket.BucketName).Set(float64(stat.ReplicaNumber))
		//example - пример вывода stats
		//TODO: подумать, как это завернуть в метрики прома. Пока думаю, что надо прать первый семпл, и отдавать его как метрику за минуту. (учитывая что семплы мы за минуту и выбираем по дефолту)
		monStats, err := couch.GetBucketStats(bucket.BucketName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		log.Println(monStats)
		for itemName, v := range monStats {
			log.Println("BucketName:", bucket.BucketName, v)
			for _, item := range v.Value {
				// TODO: нужно проверять числа вида 2.28170137e+08 и приводить их к сайзингу
				val := item[len(item)-1]
				Stats.WithLabelValues(bucket.BucketName, itemName).Set(val)
			}
		}
	}

}
