package src

import (
	"fmt"
	"net"
	"os"
	"time"
)

func handleClientUDP(conn *net.UDPConn) {
	var buf [512]byte

	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	fmt.Println("Received ", string(buf[0:]))

	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}

func checkErrorUDP(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func StartServerUDP() {
	service := ":1200"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkErrorUDP(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkErrorUDP(err)

	for {
		handleClientUDP(conn)
	}
}
