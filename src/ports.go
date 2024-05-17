package src

import (
	"net"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	openPorts = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "open_ports",
			Help: "Open ports.",
		},
		[]string{"port"},
	)
)

func init() {
	prometheus.MustRegister(openPorts)
}

func CheckOpenPorts() {
	for port := 1; port <= 65535; port++ {
		go func(port int) {
			_, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
			if err == nil {
				openPorts.WithLabelValues(strconv.Itoa(port)).Inc()
			}
		}(port)
	}
}
