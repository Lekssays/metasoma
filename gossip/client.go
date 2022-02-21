package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/Lekssays/ADeLe/gossip/proto/message"
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

	memoryProto := message.Memory{
		From: 1705405968,
		Target: 1159684051,
		Checksum: "ac9b54934b516627f9e44d2965c6e4f1bb3f8aeafff24ad7178c2ee3f07cbc21",
		Signature: "YWM5YjU0OTM0YjUxNjYyN2Y5ZTQ0ZDI5NjVjNmU0ZjFiYjNmOGFlYWZmZjI0YWQ3MTc4YzJlZTNmMDdjYmMyMQ==",
		Content: []uint32{1326285247, 1118116757, 2295028482},
	}

	data, err := proto.Marshal(&memoryProto)
	if err != nil {
		log.Println(err.Error())
	}	
	// memories := GenerateMemories()
	// memoriesBytes := string(PreparePayload(memories))

	log.Printf("Sending memories as a batch to %s:%d", address, port)
	// fmt.Fprintf(conn, memoriesBytes)
	// fmt.Fprintf(conn, data)
	conn.Write(data)

	_, err = bufio.NewReader(conn).Read(buffer)
	if err == nil {
		log.Println(string(buffer))
	} else {
		log.Println(err)
	}

	conn.Close()
}
