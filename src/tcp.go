package src

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"net"
	"os"
	"time"
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestsError, TotalRequests, portChecks, requestLatency)
}

func handleClientTCP(conn net.Conn) {
	var buf [512]byte

	start := time.Now() // Запоминаем время начала обработки запроса

	_, err := conn.Read(buf[0:])
	if err != nil {
		httpRequestsError.Inc() // Увеличиваем счетчик ошибок
		return
	}

	fmt.Println("Received ", string(buf[0:]))

	daytime := time.Now().String()
	conn.Write([]byte(daytime))
	conn.Close()

	// Увеличиваем счетчик на 1
	TotalRequests.Inc()
	httpRequestsTotal.Inc() // Увеличиваем счетчик общего количества запросов

	elapsed := time.Since(start)              // Вычисляем время обработки запроса
	requestLatency.Observe(elapsed.Seconds()) // Добавляем значение в гистограмму задержек
}

func checkErrorTCP(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func StartServerTCP() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErrorTCP(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkErrorTCP(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClientTCP(conn)
	}
}
