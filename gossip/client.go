package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
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

	memories := GenerateDummyMemories()
	peeringProofs := GenerateDummyPeeringProofs()
	messageProto := PrepareMessage(memories, peeringProofs)

	data, err := proto.Marshal(&messageProto)
	if err != nil {
		log.Println(err.Error())
	}	

	log.Printf("[*] CLIENT: Sending messsage to %s:%d", address, port)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
