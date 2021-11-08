package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

import (
	"fmt"
	"log"
	"net"
)

func main() {
	print(C.initialize(false))
	
	address := "0.0.0.0"
	port := 1337
	log.Printf("Listening on %s:%d\n", address, port)
	buffer := make([]byte, 512)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		size, remoteaddr, err := ser.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(buffer, size, remoteaddr)
		var compressedMemorySize int = int(C.compressed_memory_size())
		var numberOfEntries int = int(size / (4 + compressedMemorySize))
		fmt.Println(numberOfEntries)
		// on_memories_received(&buffer[0], &buffer[4*num_entries], num_entries)
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
