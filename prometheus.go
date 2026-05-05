package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var wafProcessLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "waf_processing_latency_ms",
	Help:    "Processing time for coraza WAF",
	Buckets: []float64{0.0001, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
})

var wafOperationsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "waf_operations_count",
	Help: "WAF processing by results",
},
	[]string{"phase", "action"},
)

var wafRuleCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "waf_rule_match_count",
		Help: "What rules the WAF has matched against",
	}, []string{
		"rule_id",
	},
)

func RegisterMetrics() {
	_ = G.Metrics.Registry.Register(wafProcessLatency)
	_ = G.Metrics.Registry.Register(wafOperationsCount)
	_ = G.Metrics.Registry.Register(wafRuleCount)
	_ = G.Metrics.Registry.Register(collectors.NewGoCollector())
	_ = G.Metrics.Registry.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
