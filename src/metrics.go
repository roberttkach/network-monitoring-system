package src

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
	)
	HttpRequestsError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_error_total",
			Help: "Total number of HTTP requests that resulted in an error.",
		},
	)
	TotalRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_requests",
		Help: "The total number of requests across all services",
	})
	OpenPorts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "open_ports",
			Help: "Open ports.",
		},
		[]string{"port"},
	)
	PortChecks = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "port_checks_total",
			Help: "Total number of port checks.",
		},
	)
	RequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Histogram of the request latency.",
			Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
		},
	)
)
