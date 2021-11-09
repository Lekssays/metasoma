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
	memoriesBytes := ""
	var owners []uint32 

	for i := 0; i < len(memories); i++ {
		memoriesBytes += string(Prepare(memories[i]))
		owners = append(owners, memories[i].Owner)
	}
	log.Println("Sent memories as a batch for", owners)
	fmt.Fprintf(conn, memoriesBytes)
	owners = nil

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
