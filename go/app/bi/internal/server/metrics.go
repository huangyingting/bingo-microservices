package server

import "github.com/prometheus/client_golang/prometheus"

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "bingo",
		Subsystem: "bi",
		Name:      "request_duration",
		Help:      "requests duration(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "bingo",
		Subsystem: "bi",
		Name:      "requests_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)
