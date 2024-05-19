package src

import (
	"fmt"
	"net"

	"os"
	"time"
)

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	start := time.Now() // Запоминаем время начала обработки запроса

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		HttpRequestsError.Inc() // Увеличиваем счетчик ошибок
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
	HttpRequestsTotal.Inc() // Увеличиваем счетчик общего количества запросов

	elapsed := time.Since(start)              // Вычисляем время обработки запроса
	RequestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек
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
