package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/Lekssays/ADeLe/gossip/proto/message"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// print(C.initialize(false))

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

	// Example of handling a packet
	// filePath := "kitsune_test.csv"
	// records := ReadCSV(filePath)
	// packets := ParsePackets(records)
	// fmt.Println(packets)
	// C.on_packet_received(C.uint(packets[0].srcIp), C.uint(packets[0].dstIp), (*C.float)(&packets[0].features[0]))
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

		log.Println(messageProto)
		go sendResponse(conn)
	}
}

func sendResponse(conn net.Conn) {
	response := fmt.Sprintf("[*] SERVER: Message Received!")
	log.Println(response)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Println(err)
	}
}

// func send_memories(conn *net.UDPConn) {
// 	peer := net.UDPAddr{
// 		Port: 1337,
// 		IP:  C.get_random_peer(),
// 	}
// 	num_entries = 10
// 	buffer := make([]byte, num_entries*(4+C.compressed_memory_size()))
// 	num_entries = C.get_memories_to_share(&buffer[0], &buffer[4*num_entries], num_entries)
// 	_, err := conn.WriteToUDP(buffer.slice(0, num_entries*(4+C.compressed_memory_size())), peer)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
