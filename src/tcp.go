package src

import (
	"fmt"
	"net"
	"os"
	"time"
)

func handleClient(conn net.Conn) {
	var buf [512]byte

	_, err := conn.Read(buf[0:])
	if err != nil {
		return
	}

	fmt.Println("Received ", string(buf[0:]))

	daytime := time.Now().String()
	conn.Write([]byte(daytime))
	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func startServer() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}
