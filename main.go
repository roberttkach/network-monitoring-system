package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"network-monitoring-system/src"
)

func main() {
	// Register metrics collectors
	go src.StartServerHTTP()
	go src.CheckOpenPorts("6060")
	go src.StartServerUDP()
	go src.StartServerDNS()
	go src.StartServerTCP()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9080", nil))
}
