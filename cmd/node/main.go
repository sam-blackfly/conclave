package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	REGISTER    = 0
	ACKNOWLEDGE = 1
)

func main() {
	// register onto the registry
	err := register("127.0.0.1:1053")
	if err != nil {
		log.Fatalf("registration error, %v", err)
	}

	// listen to incoming udp packets
	port := randomPort()
	start(port)
}

func register(addr string) error {
	p := make([]byte, 1)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("connection error, %v", err)
	}

	log.Printf("[%s] <= REGISTER\n", addr)
	_, err = fmt.Fprintf(conn, "%d", REGISTER)
	if err != nil {
		return fmt.Errorf("registration error, %v", err)
	}

	_, err = bufio.NewReader(conn).Read(p)
	if err != nil {
		return fmt.Errorf("read error, %v", err)
	}

	code, err := strconv.Atoi(string(p))
	if err != nil {
		return fmt.Errorf("response code decode error, %v", err)
	}

	if code == ACKNOWLEDGE {
		log.Printf("[%s] => ACKNOWLEDGE\n", addr)
		return nil
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("connection closing error, %v", err)
	}

	return nil
}

func start(port int) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	// listen to incoming udp packets
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer closeConnection(conn)

	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error, %v\n", err)
			continue
		}

		go process(conn, addr, buf[:n])
	}
}

func process(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	log.Printf("[%s] => %s\n", addr, buf)
	// 0 - 1: ID
	// 2: QR(1): Opcode(4)
	buf[2] |= 0x80 // Set QR bit

	_, err := conn.WriteTo(buf, addr)
	if err != nil {
		log.Printf("could not write to %v\n", addr.String())
	}
}

func closeConnection(conn *net.UDPConn) {
	err := conn.Close()
	if err != nil {
		log.Printf("connection closing error, %v\n", err)
	}
}

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	min := 30000
	max := 40000

	return rand.Intn(max-min+1) + min
}
