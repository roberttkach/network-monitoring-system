package src

import (
	"fmt"
	"net"
	"time"
)

func HandleClientTCP(conn net.Conn) {
	var buf [512]byte

	start := time.Now()

	_, err := conn.Read(buf[0:])
	if err != nil {
		HttpRequestsError.Inc()
		return
	}

	fmt.Println("Received ", string(buf[0:]))

	daytime := time.Now().Format(time.RFC3339)
	conn.Write([]byte(daytime))
	conn.Close()

	// Увеличиваем счетчик на 1
	TotalRequests.Inc()
	HttpRequestsTotal.Inc()

	elapsed := time.Since(start)
	RequestLatency.Observe(elapsed.Seconds())
}

func CheckErrorTCP(err error) error {
	if err != nil {
		return fmt.Errorf("Fatal error: %w", err)
	}
	return nil
}

func StartServerTCP() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	err = CheckErrorTCP(err)
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	err = CheckErrorTCP(err)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go HandleClientTCP(conn)
	}
}
