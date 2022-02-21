package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/Lekssays/ADeLe/gossip/proto/message"
	"io"
	"log"
	"net"
	"os"
)

func sendResponse(conn net.Conn, memoryProto message.Memory) {
	response := fmt.Sprintf("[*] SERVER: Received Memory of Target %d From %d", memoryProto.Target, memoryProto.From)
	log.Println(response)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Printf("[*] SERVER: Handling request from %s\n", conn.RemoteAddr().String())
	buffer := make([]byte, 2048)
	for {
		len, err := conn.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		memoryProto := message.Memory{}
		err = proto.Unmarshal(buffer[:len], &memoryProto)
		if err != nil {
			log.Println(err.Error())
		}

		// message := string(buffer[:len])
		// fmt.Println(message)

		go sendResponse(conn, memoryProto)
	}
}

func main() {
	address := "0.0.0.0"
	port := 1337
	addr := fmt.Sprintf("%s:%d", address, port)
	ser, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer ser.Close()

	log.Printf("Listening on %s:%d\n", address, port)
	for {
		conn, err := ser.Accept()
		if err != nil {
			log.Printf("Error accepting connection:", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
