package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

//rather than having buckets, the metric struct is used to represent
//both incoming metrics and aggregated metrics
type metric struct {
	Metric    string
	Host      string
	Timestamp string
	Service   string
	Type      string
	Value     float64
	Sampling  float64
	Tags      map[string]string
	Values    []float64 //used for histograms
	Fields    map[string]interface{}
}

type influxDBConfig struct {
	influxHost     string
	influxPort     string
	influxUsername string
	influxPassword string
	influxDatabase string
}

var (
	metrics       = make(map[string]*metric)
	in            = make(chan metric, 10000)
	out           = make(chan metric, 10000)
	flushInterval = 10 //flag.Int64("flush-interval", 10, "Flush interval")
	handlers      = make(map[string]func(metric))
	influxConfig  influxDBConfig
)

func aggregate() {
	t := time.NewTicker(time.Duration(flushInterval) * time.Second)
	for {
		select {
		case <-t.C:
			flush()

		case receivedMetric := <-in:

			if receivedMetric.Metric == "" {
				fmt.Println("Invalid Name")
				continue
			} else if receivedMetric.Timestamp == "" {
				fmt.Println("Invalid timestamp")
				continue
			} else if receivedMetric.Type == "" {
				fmt.Println("Invalid Type")
				continue
			}

			// oldMetric := metrics[receivedMetric.Metric]
			// found := false
			//
			// _, ok := metrics[receivedMetric.Metric]

			// if ok {
			// 	for _, newTag := range receivedMetric.Tags {
			// 		for _, oldTag := range oldMetric.Tags {
			// 			if newTag == oldTag {
			// 				found = true
			// 			}
			// 		}
			// 		if !found {
			// 			oldMetric.Tags = append(oldMetric.Tags, newTag)
			// 		}
			// 	}
			// }

			//if handler exists, call it.
			//ensures no exception occurs due to malformed metrics
			if handler, ok := handlers[receivedMetric.Type]; ok {
				handler(receivedMetric)
			}
		}
	}
}

//write metrics to InfluxDB
func flush() {
	if len(metrics) > 0 {
		connection := configureInfluxDB(influxConfig)
		writeInfluxDB(metrics, &connection, influxConfig)
		metrics = make(map[string]*metric)
	}

}

//http handler function, unmarshalls json encoded metric into metric struct
func receiveMetric(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var receivedMetric metric
	err := decoder.Decode(&receivedMetric)

	if err == nil {
		in <- receivedMetric
	} else {
		fmt.Println(err)
	}
}

func main() {
	viper.SetConfigName("aggregated")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("No configuration file found, exiting")
	}

	influxConfig = influxDBConfig{
		influxHost:     viper.GetString("influxHost"),
		influxPort:     viper.GetString("influxPort"),
		influxUsername: viper.GetString("influxUsername"),
		influxPassword: viper.GetString("influxPassword"),
		influxDatabase: viper.GetString("influxDatabase"),
	}

	registerHandlers()
	go aggregate()
	http.HandleFunc("/metrics", receiveMetric)
	log.Fatal(http.ListenAndServe(":8082", nil))
}