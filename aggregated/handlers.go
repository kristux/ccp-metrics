package main

import (
	"sort"
	"strconv"
)

func gaugeHandler(receivedMetric metric) {
	_, ok := metrics[receivedMetric.Metric]

	if !ok {
		metrics[receivedMetric.Metric] = &receivedMetric
		metrics[receivedMetric.Metric].Fields = make(map[string]interface{})
	} else {
		metrics[receivedMetric.Metric].Value += receivedMetric.Value
		metrics[receivedMetric.Metric].Timestamp = receivedMetric.Timestamp
		metrics[receivedMetric.Metric].Fields["gauge"] = metrics[receivedMetric.Metric].Value

	}
}

func counterHandler(receivedMetric metric) {
	_, ok := metrics[receivedMetric.Metric]

	//to avoid the metric being lost, if sampling is undefined set it to 1
	//unless the client is misbehaving, this shouldn't happen
	if receivedMetric.Sampling == 0 {
		receivedMetric.Sampling = 1
	}

	if !ok {
		metrics[receivedMetric.Metric] = &receivedMetric
		metrics[receivedMetric.Metric].Value *= 1 / receivedMetric.Sampling
		metrics[receivedMetric.Metric].Fields = make(map[string]interface{})
		metrics[receivedMetric.Metric].Fields["counter"] = metrics[receivedMetric.Metric].Value
	} else {
		metrics[receivedMetric.Metric].Value += receivedMetric.Value * (1 / receivedMetric.Sampling)
		metrics[receivedMetric.Metric].Timestamp = receivedMetric.Timestamp

	}

}

func setHandler(receivedMetric metric) {
	_, ok := metrics[receivedMetric.Metric]

	if !ok {
		metrics[receivedMetric.Metric] = &receivedMetric
		metrics[receivedMetric.Metric].Fields = make(map[string]interface{})
		//metrics[receivedMetric.Metric].Fields["items"] = make([]float64, 1)
	}

	// var set []float64
	// set = metrics[receivedMetric.Metric].Fields["items"].([]float64)
	// found := false
	//
	// for i := 0; i < len(set); i++ {
	// 	if set[i] == receivedMetric.Value {
	// 		found = true
	// 		break
	// 	}
	// }

	// if !found {
	// 	set = append(set, receivedMetric.Value)
	// }

	k := strconv.FormatFloat(float64(receivedMetric.Value), 'f', 2, 32)

	metrics[receivedMetric.Metric].Fields[k] = receivedMetric.Value
}

func histogramHandler(receivedMetric metric) {

	_, ok := metrics[receivedMetric.Metric]

	if !ok {
		metrics[receivedMetric.Metric] = &receivedMetric
		metrics[receivedMetric.Metric].Fields = make(map[string]interface{})
	}

	histogram := metrics[receivedMetric.Metric]

	histogram.Timestamp = receivedMetric.Timestamp
	histogram.Values = append(histogram.Values, receivedMetric.Value)
	sort.Float64s(histogram.Values)

	count := float64(len(histogram.Values))

	total := 0.0
	for _, x := range histogram.Values {
		total += x
	}

	//calculate stats from the values
	//go doesn't seem to have a decent stats library so none is used
	average := total / count
	median := histogram.Values[len(histogram.Values)/2]
	max := histogram.Values[int(count-1)]
	index := float64(0.95) * count
	percentile95 := histogram.Values[int(index)]

	histogram.Fields["count"] = count
	histogram.Fields["avg"] = average
	histogram.Fields["median"] = median
	histogram.Fields["max"] = max
	histogram.Fields["95percentile"] = percentile95

}

func registerHandlers() {
	handlers["gauge"] = gaugeHandler
	handlers["gauge"] = gaugeHandler
	handlers["set"] = setHandler
	handlers["counter"] = counterHandler
	handlers["histogram"] = histogramHandler
}
