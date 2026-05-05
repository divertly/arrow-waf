package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var shortlinkMetrics = prometheus.NewCounterVec( // create new counter metric. This is replacement for `prometheus.Metric` struct
	prometheus.CounterOpts{
		Name: "divertly_shortlink_hits",
		Help: "How many hits to shortlinks partitioned by domain and route",
	},
	[]string{"domain", "route", "status"},
)

var apiMetrics = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "divertly_api_http_calls",
		Help: "Requests to the API server by status",
	},
	[]string{"endpoint", "status"},
)

func RegisterMetrics() {
	_ = G.Metrics.Registry.Register(shortlinkMetrics)
	_ = G.Metrics.Registry.Register(apiMetrics)
	_ = G.Metrics.Registry.Register(collectors.NewGoCollector())
	_ = G.Metrics.Registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
