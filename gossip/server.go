package main

import (
	"fmt"
	"github.com/Lekssays/ADeLe/gossip/proto/message"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"os"
)

func sendResponse(conn net.Conn) {
	response := fmt.Sprintf("[*] SERVER: Message Received!")
	log.Println(response)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(conn net.Conn) {
	log.Println("[*] SERVER: Handling request from", conn.RemoteAddr().String())
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

		messageProto := message.Message{}
		err = proto.Unmarshal(buffer[:len], &messageProto)
		if err != nil {
			log.Println(err.Error())
		}

		go sendResponse(conn)
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
