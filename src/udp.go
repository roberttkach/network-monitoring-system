package src

import (
	"fmt"
	"net"
	"os"
	"time"
)

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	fmt.Println("Received ", string(buf[0:]))

	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func StartServer() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		handleClient(conn)
	}
}
