package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/username/gonetmon/netmon"
)

var (
	netTraffic = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_traffic_bytes",
			Help: "Network traffic.",
		},
		[]string{"device", "direction"},
	)
)

func init() {
	prometheus.MustRegister(netTraffic)
}

func main() {
	go func() {
		netInterfaces, err := netmon.GetNetInterfaces()
		if err != nil {
			log.Fatalf("Failed to get network interfaces: %v", err)
		}

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for _, netInterface := range netInterfaces {
				netTraffic.WithLabelValues(netInterface.Name, "rx").Set(float64(netInterface.RxBytes))
				netTraffic.WithLabelValues(netInterface.Name, "tx").Set(float64(netInterface.TxBytes))
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
