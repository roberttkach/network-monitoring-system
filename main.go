package main

import (
	"log"
	"net/http"
	"network-monitoring-system/src"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	go src.StartServer()
	go src.CheckOpenPorts()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
