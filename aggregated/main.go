package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type metric struct {
	Name      string
	Host      string
	Timestamp string
	Type      string
	Value     float64
	Sampling  float64
	Tags      map[string]string
}

type bucket struct {
	Name      string
	Timestamp string
	Tags      map[string]string
	Values    []float64 //used for histograms
	Fields    map[string]interface{}
}

type event struct {
	Name           string
	Text           string
	Host           string
	AggregationKey string
	Priority       string
	Timestamp      string
	AlertType      string
	Tags           map[string]string
	Fields         map[string]interface{}
}

type influxDBConfig struct {
	influxHost     string
	influxPort     string
	influxUsername string
	influxPassword string
	influxDatabase string
}

var (
	metricsIn     = make(chan metric, 10000)
	eventsIn      = make(chan event, 10000)
	flushInterval = 10 //flag.Int64("flush-interval", 10, "Flush interval")
	aggregators   = make(map[string]func(metric))
	influxConfig  influxDBConfig
	buckets       = make(map[string]*bucket)
)

func aggregate() {
	t := time.NewTicker(time.Duration(flushInterval) * time.Second)
	for {
		select {
		case <-t.C:
			flush()

		case receivedMetric := <-metricsIn:
			//if a handler exists to aggregate the metric, do so
			//otherwise ignore the metric
			if handler, ok := aggregators[receivedMetric.Type]; ok {

				if receivedMetric.Timestamp == "" {
					fmt.Println("Invalid timestamp")
					continue
				} else if receivedMetric.Type == "" {
					fmt.Println("Invalid Type")
					continue
				}

				_, ok := buckets[receivedMetric.Name]

				if !ok {
					buckets[receivedMetric.Name] = new(bucket)
					buckets[receivedMetric.Name].Name = receivedMetric.Name
					buckets[receivedMetric.Name].Fields = make(map[string]interface{})
					buckets[receivedMetric.Name].Tags = make(map[string]string)
				}

				//update tags
				//this results in the aggregated metric having the tags from the last metric
				//maybe not best, think about alternative approaches
				for k, v := range receivedMetric.Tags {
					buckets[receivedMetric.Name].Tags[k] = v
				}

				handler(receivedMetric)
			}
		}
	}
}

//write metrics to InfluxDB
func flush() {
	if len(buckets) > 0 {
		client := configureInfluxDB(influxConfig)
		writeInfluxDB(buckets, &client, influxConfig)
		buckets = make(map[string]*bucket)
	}

}

//http handler function, unmarshalls json encoded metric into metric struct
func receiveMetric(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var receivedMetric metric
	err := decoder.Decode(&receivedMetric)

	if err == nil {
		metricsIn <- receivedMetric
	} else {
		fmt.Println(err)
	}
}

func receiveEvent(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var receivedEvent event
	err := decoder.Decode(&receivedEvent)

	if err == nil {
		eventsIn <- receivedEvent
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

	registerAggregators()
	go aggregate()
	http.HandleFunc("/metrics", receiveMetric)
	http.HandleFunc("/events", receiveEvent)

	log.Fatal(http.ListenAndServe(":8082", nil))
}
