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
	buffer := make([]byte, 2048)

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

	for i := 0; i < len(memories); i++ {
		encodedMemory := Encode(memories[i])
		fmt.Fprintf(conn, string(encodedMemory))
	}

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
