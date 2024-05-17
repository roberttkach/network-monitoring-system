package src

import (
	"fmt"
	"net"
	"os"
	"time"
)

func handleClientTCP(conn net.Conn) {
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
