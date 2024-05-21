package src

import (
	"net"
	"time"
)

func CheckOpenPorts(port string) {
	start := time.Now()

	_, err := net.Dial("tcp", "localhost:"+port)
	if err == nil {
		OpenPorts.WithLabelValues(port).Inc()
	}

	elapsed := time.Since(start)
	RequestLatency.Observe(elapsed.Seconds())
	PortChecks.Inc()
	TotalRequests.Inc()
}
