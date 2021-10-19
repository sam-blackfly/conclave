package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	REGISTER    = 0
	ACKNOWLEDGE = 1
)

var connections = make(map[*net.UDPAddr]bool)

func main() {
	// listen to incoming udp packets
	start(1053)
}

func start(port int) {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// listen to incoming udp packets
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer closeConnection(conn)

	for {
		buf := make([]byte, 1)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read error, %v\n", err)
			continue
		}

		go process(conn, remoteAddr, buf[:n])
	}
}

func process(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	code, err := strconv.Atoi(string(buf))
	if err != nil {
		fmt.Printf("response code decode error, %v", err)
	}

	if code == REGISTER {
		connections[addr] = true
		log.Printf("[%s] => REGISTER\n", addr)
		log.Printf("[%s] <= ACKNOWLEDGE\n", addr)
		_, err = conn.WriteTo([]byte(strconv.Itoa(ACKNOWLEDGE)), addr)
		if err != nil {
			log.Printf("could not write to %v\n", addr.String())
		}
	}
}

func closeConnection(conn *net.UDPConn) {
	err := conn.Close()
	if err != nil {
		log.Printf("connection closing error, %v\n", err)
	}
}
