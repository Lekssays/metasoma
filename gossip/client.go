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

	addr := fmt.Sprintf("%s:%d", address, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}

	// memories := GenerateMemories()
	// memoriesBytes := string(PreparePayload(memories))

	log.Printf("Sending memories as a batch to %s:%d", address, port)
	// fmt.Fprintf(conn, memoriesBytes)
	fmt.Fprintf(conn, "something that needs to be sent")
	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
