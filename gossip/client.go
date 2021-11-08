package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	address := "0.0.0.0"
	port := 1337
	buffer := make([]byte, 512)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Println(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}

	memories := GenerateMemories()
	memoriesBytes := make([]byte, 0)
	for i := 0; i < len(memories); i++ {
		AppendToBytes(&memoriesBytes, Prepare(memories[i]))
	}
	fmt.Fprintf(conn, string(memoriesBytes))

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
