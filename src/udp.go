package src

import (
	"fmt"
	"net"

	"time"
)

func HandleClient(conn *net.UDPConn) {
	var buf [512]byte

	start := time.Now()

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		HttpRequestsError.Inc()
		CheckError(err)
		return
	}

	fmt.Println("Received ", string(buf[:n]))

	daytime := time.Now().String()
	_, err = conn.WriteToUDP([]byte(daytime), addr)
	if err != nil {
		CheckError(err)
	}

	TotalRequests.Inc()
	HttpRequestsTotal.Inc()

	elapsed := time.Since(start)
	RequestLatency.Observe(elapsed.Seconds())
}

func CheckError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Fatal error: %v", err.Error()))
	}
}
func StartServerUDP() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	CheckError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	CheckError(err)

	for {
		HandleClient(conn)
	}
}
