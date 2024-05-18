package src

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"os"
	"time"
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestsError, TotalRequests)
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	start := time.Now() // Запоминаем время начала обработки запроса

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		httpRequestsError.Inc() // Увеличиваем счетчик ошибок
		checkError(err)
		return
	}

	fmt.Println("Received ", string(buf[:n]))

	daytime := time.Now().String()
	_, err = conn.WriteToUDP([]byte(daytime), addr)
	if err != nil {
		checkError(err)
	}

	// Увеличиваем счетчик на 1
	TotalRequests.Inc()
	httpRequestsTotal.Inc() // Увеличиваем счетчик общего количества запросов

	elapsed := time.Since(start)              // Вычисляем время обработки запроса
	requestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v", err.Error())
		os.Exit(1)
	}
}

func StartServerUDP() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		handleClient(conn)
	}
}
