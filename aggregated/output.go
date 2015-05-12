package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/influxdb/influxdb/client"
)

func configureInfluxDB(config influxDBConfig) client.Client {

	influxURL, err := url.Parse(fmt.Sprintf("http://%s:%s", config.influxHost, config.influxPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *influxURL,
		Username: config.influxUsername,
		Password: config.influxPassword,
	}

	influxConnection, err := client.NewClient(conf)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return *influxConnection
}

func writeInfluxDB(buckets map[string]*bucket, influxConnection *client.Client, config influxDBConfig) {
	var (
		points      = make([]client.Point, len(buckets))
		pointsIndex = 0
	)

	for k := range buckets {
		bucket := buckets[k]
		timestamp, _ := time.Parse("YYYY-MM-DD HH:MM:SS.mmm", bucket.Timestamp)

		fmt.Println(bucket.Name)
		fmt.Println(bucket.Tags)
		fmt.Println(bucket.Fields)
		fmt.Println(timestamp)

		points[pointsIndex] = client.Point{
			Name:      bucket.Name,
			Tags:      bucket.Tags,
			Fields:    bucket.Fields,
			Timestamp: timestamp,
		}
		pointsIndex++
	}

	pointsBatch := client.BatchPoints{
		Points:          points,
		Database:        config.influxDatabase,
		RetentionPolicy: "default",
	}

	_, err := influxConnection.Write(pointsBatch)
	if err != nil {
		fmt.Println(err)
	}
}
